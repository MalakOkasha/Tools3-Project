package models

import (
	"time"

	"github.com/google/uuid" // UUID package
)

// OrderStatus defines the possible statuses of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCanceled  OrderStatus = "canceled"
)

// Order represents the structure of an order in the application
type Order struct {
	ID              uuid.UUID   `json:"id"`                     // Use uuid.UUID for ID
	UserID          uuid.UUID   `json:"user_id"`                // Use uuid.UUID for user_id
	CourierID       *uuid.UUID  `json:"courier_id,omitempty"`   // Use *uuid.UUID to handle optional courier_id
	StoreID         uuid.UUID   `json:"store_id"`               // Use uuid.UUID for store_id
	ItemIDs         []uuid.UUID `json:"item_ids"`               // Use []uuid.UUID for item_ids
	TotalPrice      float64     `json:"total_price"`            // Total price of the order
	Status          string      `json:"status"`                 // Status of the order (pending, confirmed, etc.)
	PickupLocation  string      `json:"pickup_location"`        // Pickup location
	DropOffLocation string      `json:"drop_off_location"`      // Drop off location
	PackageDetails  string      `json:"package_details"`        // Additional package details
	CreatedAt       time.Time   `json:"created_at"`             // Timestamp when the order was created
	UpdatedAt       time.Time   `json:"updated_at"`             // Timestamp when the order was last updated
	DeliveredAt     *time.Time  `json:"delivered_at,omitempty"` // Optional delivered_at timestamp
}

// AddOrderRequest defines the structure for creating a new order
type AddOrderRequest struct {
	UserID  uuid.UUID   `json:"user_id"`  // User ID for the order
	ItemIDs []uuid.UUID `json:"item_ids"` // Item IDs in the order
}

// UpdateOrderStatusRequest defines the structure for updating an order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status"` // New status for the order
}
