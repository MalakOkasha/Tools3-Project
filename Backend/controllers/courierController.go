package controllers

import (
	"PTS/models"
	"PTS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

// CourierController handles courier-related operations
type CourierController struct{}

// Register godoc
// @Summary Register a new courier
// @Description Register a new courier with details such as name, email, phone, password, location, vehicle type, and store ID. Returns a success message if registration is successful.
// @Accept json
// @Produce json
// @Param courier body models.CourierRegisterRequest true "Courier registration data"
// @Success 201 {object} map[string]string "Success response message"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /couriers/register [post]
func (ac *CourierController) CourierRegister(w http.ResponseWriter, r *http.Request) {
	var req models.CourierRegisterRequest

	// Decode the request body into the CourierRegisterRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" || req.Location == "" || req.VehicleType == "" || req.StoreId == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if the store exists
	var storeExists bool
	checkStoreQuery := "SELECT EXISTS (SELECT 1 FROM stores WHERE id = $1)"
	err := utils.DB.QueryRow(checkStoreQuery, req.StoreId).Scan(&storeExists)
	if err != nil {
		log.Println("Error checking store existence:", err)
		http.Error(w, "Store does not exist", http.StatusInternalServerError)
		return
	}
	if !storeExists {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create a new user model
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		Location:  req.Location,
		CreatedAt: time.Now(),
	}

	// Insert user into the users table
	query := "INSERT INTO users (name, email, password, phone, location, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = utils.DB.QueryRow(query, user.Name, user.Email, user.Password, user.Phone, user.Location, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	// Insert courier details
	courierQuery := "INSERT INTO couriers (user_id, vehicle_type, available, last_active_at, store_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var courierID uuid.UUID
	err = utils.DB.QueryRow(courierQuery, user.ID, req.VehicleType, true, time.Now(), req.StoreId).Scan(&courierID)
	if err != nil {
		log.Println("Error inserting courier:", err)
		http.Error(w, "Could not register courier", http.StatusInternalServerError)
		return
	}

	// Update the stores table to add the new courier's ID to couriers_ids
	updateStoreQuery := "UPDATE stores SET couriers_ids = array_append(couriers_ids, $1) WHERE id = $2"
	_, err = utils.DB.Exec(updateStoreQuery, courierID, req.StoreId)
	if err != nil {
		log.Println("Error updating store with courier ID:", err)
		http.Error(w, "Could not update store with courier information", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Courier registered successfully"})
}

// Login godoc
// @Summary Login a courier
// @Description Login a courier with email and password
// @Accept json
// @Produce json
// @Param courier body models.CourierLoginRequest true "Courier login data"
// @Success 200 {object} map[string]interface{} "Success response with JWT token and courier details"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Server error"
// @Router /couriers/login [post]
func (ac *CourierController) CourierLogin(w http.ResponseWriter, r *http.Request) {
	var req models.CourierLoginRequest

	// Decode the request body into the CourierLoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	var courier models.Courier

	// Query to retrieve all fields from both users and couriers tables
	query := `
        SELECT 
            u.id, u.name, u.email, u.password, u.phone, u.location, u.created_at,
            c.vehicle_type, c.available, c.last_active_at, c.orders
        FROM users u
        JOIN couriers c ON u.id = c.user_id
        WHERE u.email = $1
    `

	err := utils.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Location, &user.CreatedAt,
		&courier.VehicleType, &courier.Available, &courier.LastActiveAt, pq.Array(&courier.AssignedOrders),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Println("Error retrieving user:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		log.Println("Error generating JWT token:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"location":   user.Location,
			"created_at": user.CreatedAt,
		},
		"courier": map[string]interface{}{
			"vehicleType": courier.VehicleType,
			"available":   courier.Available,
			"lastActive":  courier.LastActiveAt,
			"orders":      courier.AssignedOrders,
			"store_id":    courier.StoreId,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
