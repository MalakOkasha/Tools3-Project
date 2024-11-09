package controllers

import (
	"PTS/models"
	"PTS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type OrderController struct{}

// AddOrder godoc
// @Summary Add a new order
// @Description Create a new order with user ID, item IDs, and drop-off location. Store ID and courier ID are determined automatically.
// @Accept json
// @Produce json
// @Param order body models.AddOrderRequest true "Order data"
// @Success 201 {object} map[string]string "Order added successfully"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/add [post]
// @Tag Orders
func (oc *OrderController) AddOrder(w http.ResponseWriter, r *http.Request) {
	var req models.AddOrderRequest

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request fields
	if req.UserID == uuid.Nil || len(req.ItemIDs) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Start a database transaction
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Defer a rollback in case of any error
	defer tx.Rollback()

	// Retrieve the Store ID from the first item's store_id
	var storeID uuid.UUID
	storeQuery := "SELECT store_id FROM items WHERE id = $1"
	err = tx.QueryRow(storeQuery, req.ItemIDs[0]).Scan(&storeID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		log.Println("Error fetching store ID from item:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Retrieve the PickupLocation from the store's location
	var pickupLocation string
	locationQuery := "SELECT location FROM stores WHERE id = $1"
	err = tx.QueryRow(locationQuery, storeID).Scan(&pickupLocation)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Store not found", http.StatusNotFound)
			return
		}
		log.Println("Error fetching store location:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Retrieve the User's location as DropOffLocation
	var userLocation string
	userQuery := "SELECT location FROM users WHERE id = $1"
	err = tx.QueryRow(userQuery, req.UserID).Scan(&userLocation)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Println("Error fetching user location:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Calculate the total price by summing up the prices of items in the order
	var totalPrice float64
	for _, itemID := range req.ItemIDs {
		var itemPrice float64
		itemQuery := "SELECT price FROM items WHERE id = $1"
		err := tx.QueryRow(itemQuery, itemID).Scan(&itemPrice)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			log.Println("Error fetching item price:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		totalPrice += itemPrice
	}

	// Select a random available courier for the store, or leave courier_id NULL if no available couriers
	var courierID *uuid.UUID
	courierQuery := "SELECT id FROM couriers WHERE store_id = $1 AND available = true ORDER BY RANDOM() LIMIT 1"
	err = tx.QueryRow(courierQuery, storeID).Scan(&courierID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error selecting courier:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// If a courier is selected, mark them as unavailable
	if courierID != nil {
		updateCourierAvailabilityQuery := "UPDATE couriers SET available = false WHERE id = $1"
		_, err := tx.Exec(updateCourierAvailabilityQuery, *courierID)
		if err != nil {
			log.Println("Error updating courier availability:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	}

	// Create the order object
	order := models.Order{
		UserID:          req.UserID,
		CourierID:       courierID,
		StoreID:         storeID,
		ItemIDs:         req.ItemIDs,
		TotalPrice:      totalPrice,
		Status:          string(models.OrderStatusPending),
		PickupLocation:  pickupLocation,
		DropOffLocation: userLocation,
		PackageDetails:  "New Customer", // Example; can be changed as needed
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Insert the order into the orders table
	insertOrderQuery := `INSERT INTO orders (user_id, courier_id, store_id, item_ids, total_price, status, pickup_location, drop_off_location, package_details, created_at, updated_at)
						 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	var orderID uuid.UUID
	err = tx.QueryRow(insertOrderQuery, order.UserID, order.CourierID, order.StoreID, pq.Array(order.ItemIDs), order.TotalPrice, order.Status, order.PickupLocation, order.DropOffLocation, order.PackageDetails, order.CreatedAt, order.UpdatedAt).Scan(&orderID)
	if err != nil {
		log.Println("Error inserting order:", err)
		http.Error(w, "Could not add order", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Order added successfully", "order_id": orderID})
}

// DeleteOrder godoc
// @Summary Delete an order by ID
// @Description Delete an order by its ID
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order deleted successfully"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/delete/{id} [delete]
// @Tag Orders
func (oc *OrderController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM orders WHERE id = $1"
	result, err := utils.DB.Exec(query, orderID)
	if err != nil {
		log.Println("Error deleting order:", err)
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted successfully"})
}

// ListOrdersByUserID godoc
// @Summary List all orders for a user
// @Description Get a list of all orders for a specific user
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Order "List of orders"
// @Failure 400 {object} map[string]string "User ID is required"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/list/user/{user_id} [get]
// @Tag Orders
func (oc *OrderController) ListOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query("SELECT id, user_id, courier_id, store_id, item_ids, status, pickup_location, drop_off_location, package_details, created_at, updated_at, delivered_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.CourierID, &order.StoreID, pq.Array(&order.ItemIDs), &order.Status, &order.PickupLocation, &order.DropOffLocation, &order.PackageDetails, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		http.Error(w, "No orders found for the user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrderDetails godoc
// @Summary Get detailed order information
// @Description Retrieve all details about an order, including user, store, items, courier, and status.
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Order details retrieved successfully"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/{order_id} [get]
// @Tag Orders
func (oc *OrderController) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	// Retrieve the order ID from the URL
	orderID := mux.Vars(r)["order_id"]

	// SQL query to get order details including user, store, courier, and items
	query := `
		SELECT
			o.id, o.user_id, o.courier_id, o.store_id, o.item_ids,
			o.total_price, o.status, o.pickup_location, o.drop_off_location,
			o.package_details, o.created_at, o.updated_at, o.delivered_at,
			u.id, u.name, u.email, u.phone, u.location, u.created_at,
			s.id, s.name, s.location,
			c.id  -- Adjusted to select only the existing columns
		FROM
			orders o
		LEFT JOIN users u ON o.user_id = u.id
		LEFT JOIN stores s ON o.store_id = s.id
		LEFT JOIN couriers c ON o.courier_id = c.id
		WHERE o.id = $1
	`

	// Execute the query
	row := utils.DB.QueryRow(query, orderID)

	// Initialize variables to store query results
	var order models.Order
	var user models.User
	var store models.Store
	var courier models.Courier

	// Scan the result into variables, using pq.Array to scan the item_ids as []uuid.UUID
	err := row.Scan(
		&order.ID, &order.UserID, &order.CourierID, &order.StoreID, pq.Array(&order.ItemIDs),
		&order.TotalPrice, &order.Status, &order.PickupLocation, &order.DropOffLocation,
		&order.PackageDetails, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt,
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Location, &user.CreatedAt,
		&store.ID, &store.Name, &store.Location,
		&courier.ID, // Adjusted for the missing `phone` column
	)

	// Error handling
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving order details: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Fetch item details separately
	itemsQuery := `
		SELECT id, name, price
		FROM items
		WHERE id = ANY($1)
	`
	itemRows, err := utils.DB.Query(itemsQuery, pq.Array(order.ItemIDs))
	if err != nil {
		log.Printf("Error retrieving items: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer itemRows.Close()

	// Collect item details
	var items []map[string]interface{}
	for itemRows.Next() {
		var item models.Item
		if err := itemRows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Error scanning item: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		items = append(items, map[string]interface{}{
			"id":    item.ID,
			"name":  item.Name,
			"price": item.Price,
		})
	}

	// Construct the response
	orderDetails := map[string]interface{}{
		"order": map[string]interface{}{
			"id":                order.ID,
			"user_id":           order.UserID,
			"courier_id":        order.CourierID,
			"store_id":          order.StoreID,
			"item_ids":          order.ItemIDs,
			"total_price":       order.TotalPrice,
			"status":            order.Status,
			"pickup_location":   order.PickupLocation,
			"drop_off_location": order.DropOffLocation,
			"package_details":   order.PackageDetails,
			"created_at":        order.CreatedAt,
			"updated_at":        order.UpdatedAt,
			"delivered_at":      order.DeliveredAt,
		},
		"user": map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"location":   user.Location,
			"created_at": user.CreatedAt,
		},
		"store": map[string]interface{}{
			"id":       store.ID,
			"name":     store.Name,
			"location": store.Location,
		},
		"courier": map[string]interface{}{
			"id": courier.ID, // Adjusted for the missing `phone` column
		},
		"items": items,
	}

	// Respond with the order details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderDetails)
}

// ListOrdersByStoreID godoc
// @Summary List all orders for a store
// @Description Get a list of all orders for a specific store
// @Accept json
// @Produce json
// @Param store_id path string true "Store ID"
// @Success 200 {array} models.Order "List of orders"
// @Failure 400 {object} map[string]string "Store ID is required"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/list/store/{store_id} [get]
// @Tag Orders
func (oc *OrderController) ListOrdersByStoreID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["store_id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query("SELECT id, user_id, courier_id, store_id, item_ids, status, pickup_location, drop_off_location, package_details, created_at, updated_at, delivered_at FROM orders WHERE store_id = $1", storeID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.CourierID, &order.StoreID, pq.Array(&order.ItemIDs), &order.Status, &order.PickupLocation, &order.DropOffLocation, &order.PackageDetails, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		http.Error(w, "No orders found for the store", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// ListOrdersByCourierID godoc
// @Summary List all orders assigned to a courier
// @Description Get a list of all orders for a specific courier
// @Accept json
// @Produce json
// @Param courier_id path string true "Courier ID"
// @Success 200 {array} models.Order "List of orders"
// @Failure 400 {object} map[string]string "Courier ID is required"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/list/courier/{courier_id} [get]
// @Tag Orders
func (oc *OrderController) ListOrdersByCourierID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courierID := vars["courier_id"]

	if courierID == "" {
		http.Error(w, "Courier ID is required", http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query("SELECT id, user_id, courier_id, store_id, item_ids, status, pickup_location, drop_off_location, package_details, created_at, updated_at, delivered_at FROM orders WHERE courier_id = $1", courierID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.CourierID, &order.StoreID, pq.Array(&order.ItemIDs), &order.Status, &order.PickupLocation, &order.DropOffLocation, &order.PackageDetails, &order.CreatedAt, &order.UpdatedAt, &order.DeliveredAt); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		http.Error(w, "No orders found for the courier", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus godoc
// @Summary Update the status of an order
// @Description Update the status of an order by its ID. If the status is "delivered", sets `delivered_at` to the current time. If the status is "canceled", sets the courier's availability to true.
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "New Status"
// @Success 200 {object} map[string]string "Order status updated successfully"
// @Failure 400 {object} map[string]string "Order ID and status are required"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/update/{id} [patch]
// @Tag Orders
func (oc *OrderController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	// Validate the provided order ID
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Parse and validate the request body
	var req models.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// Trim any extra whitespace from the status value
	req.Status = strings.TrimSpace(req.Status)

	// Log the status for debugging purposes
	log.Printf("Received status: '%s'", req.Status)

	// Validate that the status is one of the defined OrderStatus values
	switch req.Status {
	case string(models.OrderStatusPending), string(models.OrderStatusConfirmed), string(models.OrderStatusShipped), string(models.OrderStatusDelivered), string(models.OrderStatusCanceled):
		// Valid status, proceed
	default:
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Initialize the update query and parameters
	query := "UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3"
	params := []interface{}{req.Status, time.Now(), orderID}

	// Additional updates based on status
	if req.Status == string(models.OrderStatusDelivered) {
		// If the status is "delivered", update `delivered_at` timestamp
		query = "UPDATE orders SET status = $1, updated_at = $2, delivered_at = $3 WHERE id = $4"
		params = []interface{}{req.Status, time.Now(), time.Now(), orderID}
	} else if req.Status == string(models.OrderStatusCanceled) {
		// If the status is "canceled", set the assigned courier's availability to true
		courierUpdateQuery := `
            UPDATE couriers 
            SET available = true 
            WHERE id = (SELECT courier_id FROM orders WHERE id = $1)
        `
		_, err := utils.DB.Exec(courierUpdateQuery, orderID)
		if err != nil {
			log.Println("Error updating courier availability:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	}

	// Execute the order update query
	result, err := utils.DB.Exec(query, params...)
	if err != nil {
		log.Println("Error updating order status:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Check if the order exists (i.e., rowsAffected > 0)
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated successfully"})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order by its ID and set the assigned courier's availability to true
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order canceled successfully"
// @Failure 400 {object} map[string]string "Order ID is required"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /orders/cancel/{id} [patch]
// @Tag Orders
func (oc *OrderController) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	// Validate the provided order ID
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Update the order status to "canceled" and set the courier's availability to true
	query := `
        UPDATE orders 
        SET status = 'canceled', updated_at = $1 
        WHERE id = $2 
        RETURNING courier_id
    `
	var courierID uuid.UUID
	err := utils.DB.QueryRow(query, time.Now(), orderID).Scan(&courierID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Println("Error updating order status:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Set the courier's availability to true
	courierUpdateQuery := "UPDATE couriers SET available = true WHERE id = $1"
	_, err = utils.DB.Exec(courierUpdateQuery, courierID)
	if err != nil {
		log.Println("Error updating courier availability:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order canceled successfully"})
}
