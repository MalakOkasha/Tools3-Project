package models

import "time"

type Item struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StoreID     string    `json:"store_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CoverLink   string    `json:"cover_link"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AddItemRequest represents the structure for the Add Item request
type AddItemRequest struct {
	UserID      string   `json:"user_id"`
	StoreID     string   `json:"store_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
	CoverLink   string   `json:"cover_link"`
	Images      []string `json:"images"`
}
