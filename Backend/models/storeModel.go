package models

import "time"

// Store represents a store record in the database
type Store struct {
	ID          string
	Name        string
	Location    string
	OwnerID     string
	AdminsIDs   []string
	CouriersIDs []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Store represents a store record in the database
type StoreRegisterRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	OwnerID  string `json:"owner_id"`
}

// AddAdminRequest represents the request body for adding an admin to a store
type AddAdminRequest struct {
	AdminID string `json:"adminId"`
}

// RemoveAdminRequest represents the request body for removing an admin from a store
type RemoveAdminRequest struct {
	AdminID string `json:"adminId"`
}

// AddCourierRequest represents the request body for adding a courier to a store
type AddCourierRequest struct {
	CourierID string `json:"courierId"`
}

// RemoveCourierRequest represents the request body for removing a courier from a store
type RemoveCourierRequest struct {
	CourierID string `json:"courierId"`
}
