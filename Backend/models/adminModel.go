package models

type Admin struct {
	User
	StoreId string
}

// RegisterRequest represents the structure for the registration request
type AdminRegisterRequest struct {
	Email    string `json:"email"`
	Location string `json:"location"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	StoreId  string `json:"store_id"`
}

// LoginRequest represents the structure for the login request
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
