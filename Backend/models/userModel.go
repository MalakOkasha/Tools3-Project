package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Phone     string
	CreatedAt time.Time
}

// Method to display user information
func (u *User) DisplayInfo() {
	fmt.Printf("User: %s, Email: %s, Phone: %s\n", u.Name, u.Email, u.Phone)
}

// Method to update user information
func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}
