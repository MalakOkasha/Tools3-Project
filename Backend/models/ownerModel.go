package models

type Owner struct {
	User
	StoreName     string
	StoreLocation string
	StoreId       string
}

// RegisterRequest represents the structure for the registration request
type OwnerRegisterRequest struct {
	Email         string `json:"email"`
	Location      string `json:"location"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	StoreName     string `json:"store_name"`
	StoreLocation string `json:"store_location"`
}

// LoginRequest represents the structure for the login request
type OwnerLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
