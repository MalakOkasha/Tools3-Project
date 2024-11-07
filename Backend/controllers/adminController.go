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

type AdminController struct{}

// Register godoc
// @Summary Register a new admin
// @Description Register a new admin with details such as name, email, phone, password, location, and store ID. Returns a success message if registration is successful.
// @Accept json
// @Produce json
// @Param admin body models.AdminRegisterRequest true "Admin registration data"
// @Success 201 {object} map[string]string "Success response message"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /admins/register [post]
func (ac *AdminController) AdminRegister(w http.ResponseWriter, r *http.Request) {
	var req models.AdminRegisterRequest

	// Decode the request body into the AdminRegisterRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" || req.Location == "" || req.StoreId == "" {
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
	userQuery := "INSERT INTO users (name, email, password, phone, location, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = utils.DB.QueryRow(userQuery, user.Name, user.Email, user.Password, user.Phone, user.Location, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	// Insert admin details
	adminQuery := "INSERT INTO admins (user_id, store_id) VALUES ($1, $2) RETURNING id"
	var adminID uuid.UUID
	err = utils.DB.QueryRow(adminQuery, user.ID, req.StoreId).Scan(&adminID)
	if err != nil {
		log.Println("Error inserting admin:", err)
		http.Error(w, "Could not register admin", http.StatusInternalServerError)
		return
	}

	// Update the stores table to add the new admin's ID to admin_ids
	updateStoreQuery := "UPDATE stores SET admins_ids = array_append(admin_ids, $1) WHERE id = $2"
	_, err = utils.DB.Exec(updateStoreQuery, adminID, req.StoreId)
	if err != nil {
		log.Println("Error updating store with admin ID:", err)
		http.Error(w, "Could not update store with admin information", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin registered successfully"})
}

// AdminLogin godoc
// @Summary Login an admin
// @Description Login an admin with email and password
// @Accept json
// @Produce json
// @Param admin body models.AdminLoginRequest true "Admin login data"
// @Success 200 {object} map[string]interface{} "Success response with JWT token and admin details"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Server error"
// @Router /admins/login [post]
func (ac *AdminController) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req models.AdminLoginRequest

	// Decode the request body into the AdminLoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	var admin models.Admin

	// Query to retrieve all fields from both users and admins tables
	query := `
        SELECT 
            u.id, u.name, u.email, u.password, u.phone, u.location, u.created_at,
            a.store_id
        FROM users u
        JOIN admins a ON u.id = a.user_id
        WHERE u.email = $1
    `

	err := utils.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Location, &user.CreatedAt,
		&admin.StoreId,
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
		"admin": map[string]interface{}{
			"store_id": admin.StoreId,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
