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

	"golang.org/x/crypto/bcrypt"
)

type OwnerController struct{}

// Register godoc
// @Summary Register a new owner
// @Description Register a new owner with details such as name, email, phone, password, location, and store details. Returns a success message if registration is successful.
// @Accept json
// @Produce json
// @Param owner body models.OwnerRegisterRequest true "Owner registration data"
// @Success 201 {object} map[string]string "Success response message"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /owners/register [post]
func (oc *OwnerController) OwnerRegister(w http.ResponseWriter, r *http.Request) {
	var req models.OwnerRegisterRequest

	// Decode the request body into the OwnerRegisterRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" || req.Location == "" || req.StoreName == "" || req.StoreLocation == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
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
	userQuery := "INSERT INTO users (name, email, password, phone, location, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = utils.DB.QueryRow(userQuery, user.Name, user.Email, user.Password, user.Phone, user.Location, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	// Insert store linked to the owner
	storeQuery := "INSERT INTO stores (name, location, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var storeID uuid.UUID
	err = utils.DB.QueryRow(storeQuery, req.StoreName, req.StoreLocation, user.ID, time.Now(), time.Now()).Scan(&storeID)
	if err != nil {
		log.Println("Error inserting store:", err)
		http.Error(w, "Could not register store", http.StatusInternalServerError)
		return
	}

	// Insert Owner details
	ownerQuery := "INSERT INTO owners (user_id, store_name, store_location, store_id) VALUES ($1, $2, $3, $4) RETURNING id"
	var ownerID uuid.UUID
	err = utils.DB.QueryRow(ownerQuery, user.ID, req.StoreName, req.StoreLocation, storeID).Scan(&ownerID)
	if err != nil {
		log.Println("Error inserting owner:", err)
		http.Error(w, "Could not register owner", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Owner registered successfully"})
}

// OwnerLogin godoc
// @Summary Login an owner
// @Description Login an owner with email and password
// @Accept json
// @Produce json
// @Param owner body models.OwnerLoginRequest true "Owner login data"
// @Success 200 {object} map[string]interface{} "Success response with JWT token and owner details"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Server error"
// @Router /owners/login [post]
func (oc *OwnerController) OwnerLogin(w http.ResponseWriter, r *http.Request) {
	var req models.OwnerLoginRequest

	// Decode the request body into the OwnerLoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	var owner models.Owner

	// Query to retrieve all fields from both users and owners tables
	query := `
        SELECT 
            u.id, u.name, u.email, u.password, u.phone, u.location, u.created_at,
            o.store_name, o.store_location, o.store_id
        FROM users u
        JOIN owners o ON u.id = o.user_id
        WHERE u.email = $1
    `

	err := utils.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Location, &user.CreatedAt,
		&owner.StoreName, &owner.StoreLocation, &owner.StoreId,
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
		"owner": map[string]interface{}{
			"store_id":       owner.StoreId,
			"store_name":     owner.StoreName,
			"store_location": owner.StoreLocation,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
