package models

import (
	"time"
)

type Courier struct {
	User
	VehicleType    string
	AssignedOrders []string
	Available      bool
	LastActiveAt   time.Time
	StoreId        string
}

// CourierRegisterRequest represents the structure for the courier registration request
type CourierRegisterRequest struct {
	Email       string `json:"email"`
	Location    string `json:"location"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	VehicleType string `json:"vehicle_type"`
	StoreId     string `json:"store_id"`
}

// CourierLoginRequest represents the structure for the courier login request
type CourierLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
