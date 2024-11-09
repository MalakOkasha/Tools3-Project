package controllers

import (
	"PTS/models"
	"PTS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type ItemController struct{}

// AddItem godoc
// @Summary Add a new item
// @Description Add a new item with name, description, price, stock, category, cover link, and images
// @Accept json
// @Produce json
// @Param item body models.AddItemRequest true "Item data"
// @Success 201 {object} map[string]string "Item added successfully"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /items/add [post]
// @Tag Items
func (ic *ItemController) AddItem(w http.ResponseWriter, r *http.Request) {
	var req models.AddItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Price <= 0 || req.StoreID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	item := models.Item{
		UserID:      req.UserID,
		StoreID:     req.StoreID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Stock:       req.Stock,
		CoverLink:   req.CoverLink,
		Images:      req.Images,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO items (user_id, store_id, name, description, price, stock, category, cover_link, images, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := utils.DB.Exec(query, item.UserID, item.StoreID, item.Name, item.Description, item.Price, item.Stock, item.Category, item.CoverLink, pq.Array(item.Images), item.CreatedAt, item.UpdatedAt)
	if err != nil {
		log.Println("Error inserting item:", err)
		http.Error(w, "Could not add item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added successfully"})
}

// DeleteItem godoc
// @Summary Delete an item by ID
// @Description Delete an item by its ID
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]string "Item deleted successfully"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /items/delete/{id} [delete]
// @Tag Items
func (ic *ItemController) DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from the URL path parameter
	vars := mux.Vars(r)
	itemID := vars["id"]

	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	// SQL query to delete the item by ID
	query := "DELETE FROM items WHERE id = $1"
	result, err := utils.DB.Exec(query, itemID)
	if err != nil {
		log.Println("Error deleting item:", err)
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected (i.e., if the item exists)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}

// ListItemsByStoreID godoc
// @Summary List all items for a store
// @Description Get a list of all items for a specific store
// @Accept json
// @Produce json
// @Param store_id path string true "Store ID"
// @Success 200 {array} models.Item "List of items"
// @Failure 400 {object} map[string]string "Store ID is required"
// @Failure 500 {object} map[string]string "Server error"
// @Router /items/list/{store_id} [get]
// @Tag Items
func (ic *ItemController) ListItemsByStoreID(w http.ResponseWriter, r *http.Request) {
	// Extract store ID from the URL path parameter
	vars := mux.Vars(r)
	storeID := vars["store_id"]

	// Check if the store_id is provided
	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	// Log the store ID for debugging
	log.Println("Received request to list items for store ID:", storeID)

	// Query to get all items from the store
	rows, err := utils.DB.Query("SELECT id, user_id, store_id, name, description, price, stock, category, cover_link, images, created_at, updated_at FROM items WHERE store_id = $1", storeID)
	if err != nil {
		log.Println("Error executing query to retrieve items:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the items
	var items []models.Item
	for rows.Next() {
		var item models.Item
		// Ensure images is scanned properly as an array of strings
		if err := rows.Scan(&item.ID, &item.UserID, &item.StoreID, &item.Name, &item.Description, &item.Price, &item.Stock, &item.Category, &item.CoverLink, pq.Array(&item.Images), &item.CreatedAt, &item.UpdatedAt); err != nil {
			log.Println("Error scanning item:", err)
			continue
		}
		items = append(items, item)
	}

	// Check if there are no items found
	if len(items) == 0 {
		log.Println("No items found for store ID:", storeID)
		http.Error(w, "No items found for the store", http.StatusNotFound)
		return
	}

	// Respond with the list of items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// GetItemByID godoc
// @Summary Get item by ID
// @Description Get details of an item by its ID
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} models.Item "Item details"
// @Failure 400 {object} map[string]string "Item ID is required"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /items/get/{id} [get]
// @Tag Items
func (ic *ItemController) GetItemByID(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from the URL path parameter
	vars := mux.Vars(r)
	itemID := vars["id"]

	// Check if the item ID is provided
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	// Log the item ID for debugging
	log.Println("Received request to get item with ID:", itemID)

	// Define query for fetching item by ID
	query := "SELECT id, user_id, store_id, name, description, price, stock, category, cover_link, images, created_at, updated_at FROM items WHERE id = $1"
	var item models.Item

	// Execute query and scan the result into the item struct
	err := utils.DB.QueryRow(query, itemID).Scan(&item.ID, &item.UserID, &item.StoreID, &item.Name, &item.Description, &item.Price, &item.Stock, &item.Category, &item.CoverLink, pq.Array(&item.Images), &item.CreatedAt, &item.UpdatedAt)

	// Check for errors
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Item not found with ID:", itemID)
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		log.Println("Error retrieving item with ID:", itemID, "Error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Respond with the item details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
