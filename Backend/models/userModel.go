package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Phone     string
	Location  string
	CreatedAt time.Time
}

// RegisterRequest represents the structure for the registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Location string `json:"location"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// LoginRequest represents the structure for the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Method to display user information
func (u *User) DisplayInfo() {
	fmt.Printf("User: %s, Email: %s, Phone: %s, Location: %s\n", u.Name, u.Email, u.Phone, u.Location)
}

// Method to update user information
func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}
