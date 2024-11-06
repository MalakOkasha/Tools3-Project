package controllers

import (
	"PTS/models"
	"PTS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

// RegisterRequest represents the structure for the registration request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginRequest represents the structure for the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, phone, and password. Returns a success message if registration is successful.
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]string "Success response message"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /register [post]
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Decode the request body into the RegisterRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
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
		Password:  string(hashedPassword), // Store the hashed password
		Phone:     req.Phone,
		CreatedAt: time.Now(),
	}

	// Insert the user into the database
	query := "INSERT INTO users (name, email, password, phone, created_at) VALUES ($1, $2, $3, $4, $5)"
	if _, err = utils.DB.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.CreatedAt); err != nil {
		log.Println("Error inserting user into the database:", err)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Accept json
// @Produce json
// @Param user body LoginRequest true "User login data"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /login [post]
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Decode the request body into the LoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User

	// Check if the user exists in the database
	query := "SELECT id, name, email, password FROM users WHERE email = $1"
	if err := utils.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
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

	// Return the JWT token and user data as a response
	responseData := map[string]interface{}{
		"token": token,
		"user":  user, // Include user data in the response
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
