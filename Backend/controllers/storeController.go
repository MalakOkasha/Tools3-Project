package controllers

import (
	"PTS/models"
	"PTS/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type StoreController struct{}

// AddStore godoc
// @Summary Add a new store
// @Description Create a new store with name, location, and owner ID
// @Accept json
// @Produce json
// @Param store body models.StoreRegisterRequest true "Store data"
// @Success 201 {object} map[string]string "Store added successfully"
// @Failure 400 {object} map[string]string "Missing required fields or invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/add [post]
// @Tag Stores
func (sc *StoreController) AddStore(w http.ResponseWriter, r *http.Request) {
	var req models.StoreRegisterRequest

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request fields
	if req.Name == "" || req.OwnerID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Initialize AdminsIDs and CouriersIDs as empty arrays
	adminsIDs := []string{}
	couriersIDs := []string{}

	// Insert the store into the database
	insertQuery := `
		INSERT INTO stores (name, location, owner_id, admins_ids, couriers_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	var storeID string
	err := utils.DB.QueryRow(
		insertQuery,
		req.Name, req.Location, req.OwnerID,
		pq.Array(adminsIDs), pq.Array(couriersIDs),
		time.Now(), time.Now(),
	).Scan(&storeID)

	if err != nil {
		log.Println("Error inserting store:", err)
		http.Error(w, "Could not add store", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Store added successfully",
		"store_id": storeID,
	})
}

// GetStoreByID godoc
// @Summary Get a store by its ID
// @Description Get details of a store by its ID
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} models.Store "Store details"
// @Failure 400 {object} map[string]string "Store ID is required"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/get/{id} [get]
// @Tag Stores
func (sc *StoreController) GetStoreByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, location, owner_id, admins_ids, couriers_ids, created_at, updated_at FROM stores WHERE id = $1`
	var store models.Store

	err := utils.DB.QueryRow(query, storeID).Scan(&store.ID, &store.Name, &store.Location, &store.OwnerID, pq.Array(&store.AdminsIDs), pq.Array(&store.CouriersIDs), &store.CreatedAt, &store.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

// ListStores godoc
// @Summary List all stores
// @Description Get a list of all stores
// @Accept json
// @Produce json
// @Success 200 {array} models.Store "List of stores"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/list [get]
// @Tag Stores
func (sc *StoreController) ListStores(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.DB.Query("SELECT id, name, location, owner_id, admins_ids, couriers_ids, created_at, updated_at FROM stores")
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stores []models.Store
	for rows.Next() {
		var store models.Store
		if err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.OwnerID, pq.Array(&store.AdminsIDs), pq.Array(&store.CouriersIDs), &store.CreatedAt, &store.UpdatedAt); err != nil {
			continue
		}
		stores = append(stores, store)
	}

	if len(stores) == 0 {
		http.Error(w, "No stores found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

// DeleteStore godoc
// @Summary Delete a store by ID
// @Description Delete a store by its ID
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} map[string]string "Store deleted successfully"
// @Failure 400 {object} map[string]string "Invalid store ID"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/delete/{id} [delete]
// @Tag Stores
func (sc *StoreController) DeleteStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	deleteQuery := `DELETE FROM stores WHERE id = $1`
	result, err := utils.DB.Exec(deleteQuery, storeID)
	if err != nil {
		log.Println("Error deleting store:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Store deleted successfully"})
}

// AddAdminToStore godoc
// @Summary Add an admin to a store
// @Description Add an admin by user ID to a store if the admin exists in the admins table
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param request body models.AddAdminRequest true "Admin User ID"
// @Success 200 {object} map[string]string "Admin added successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 404 {object} map[string]string "Admin not found or is not an admin user"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/add-admin/{id} [patch]
// @Tag Stores
func (sc *StoreController) AddAdminToStore(w http.ResponseWriter, r *http.Request) {
	storeID := mux.Vars(r)["id"]

	// Validate Store ID
	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req models.AddAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.AdminID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	adminUUID, err := uuid.Parse(req.AdminID)
	if err != nil {
		http.Error(w, "Invalid Admin ID format", http.StatusBadRequest)
		return
	}

	// Check if admin exists in the admins table
	var adminExists bool
	checkAdminQuery := "SELECT EXISTS (SELECT 1 FROM admins WHERE id = $1)"
	err = utils.DB.QueryRow(checkAdminQuery, adminUUID).Scan(&adminExists)
	if err != nil {
		log.Println("Error checking admin existence:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if !adminExists {
		http.Error(w, "Admin not found or is not an admin user", http.StatusNotFound)
		return
	}

	// Check if the admin is already in the store's admins_ids array
	checkStoreAdminQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM stores 
			WHERE id = $1 AND $2 = ANY(admins_ids)
		)
	`
	var alreadyAdded bool
	err = utils.DB.QueryRow(checkStoreAdminQuery, storeID, adminUUID.String()).Scan(&alreadyAdded)
	if err != nil {
		log.Println("Error checking store admin existence:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if alreadyAdded {
		http.Error(w, "Admin already exists in the store", http.StatusBadRequest)
		return
	}

	// Add the admin to the store's admins_ids array
	updateQuery := `
		UPDATE stores 
		SET admins_ids = ARRAY_APPEND(admins_ids, $1), updated_at = NOW()
		WHERE id = $2
	`
	_, err = utils.DB.Exec(updateQuery, adminUUID.String(), storeID)
	if err != nil {
		log.Println("Error adding admin to store:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin added successfully"})
}

// RemoveAdminFromStore godoc
// @Summary Remove an admin from a store
// @Description Remove an admin by user ID from a store
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param request body models.RemoveAdminRequest true "Admin User ID"
// @Success 200 {object} map[string]string "Admin removed successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/remove-admin/{id} [patch]
// @Tag Stores
func (sc *StoreController) RemoveAdminFromStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	// Parse and validate the request body
	var req models.RemoveAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.AdminID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateQuery := `UPDATE stores SET admins_ids = array_remove(admins_ids, $1), updated_at = $2 WHERE id = $3`
	_, err := utils.DB.Exec(updateQuery, req.AdminID, time.Now(), storeID)
	if err != nil {
		log.Println("Error removing admin from store:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin removed successfully"})
}

// AddCourierToStore godoc
// @Summary Add a courier to a store
// @Description Add a courier by user ID to a store if the courier exists in the couriers table
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param request body models.AddCourierRequest true "Courier User ID"
// @Success 200 {object} map[string]string "Courier added successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 404 {object} map[string]string "Courier not found or is not a courier user"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/add-courier/{id} [patch]
// @Tag Stores
func (sc *StoreController) AddCourierToStore(w http.ResponseWriter, r *http.Request) {
	storeID := mux.Vars(r)["id"]

	// Validate Store ID
	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req models.AddCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CourierID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	courierUUID, err := uuid.Parse(req.CourierID)
	if err != nil {
		http.Error(w, "Invalid Courier ID format", http.StatusBadRequest)
		return
	}

	// Check if courier exists in the couriers table
	var courierExists bool
	checkCourierQuery := "SELECT EXISTS (SELECT 1 FROM couriers WHERE id = $1)"
	err = utils.DB.QueryRow(checkCourierQuery, courierUUID).Scan(&courierExists)
	if err != nil {
		log.Println("Error checking courier existence:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if !courierExists {
		http.Error(w, "Courier not found or is not a courier user", http.StatusNotFound)
		return
	}

	// Check if the courier is already in the store's couriers_ids array
	checkStoreCourierQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM stores 
			WHERE id = $1 AND $2 = ANY(couriers_ids)
		)
	`
	var alreadyAdded bool
	err = utils.DB.QueryRow(checkStoreCourierQuery, storeID, courierUUID.String()).Scan(&alreadyAdded)
	if err != nil {
		log.Println("Error checking store courier existence:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if alreadyAdded {
		http.Error(w, "Courier already exists in the store", http.StatusBadRequest)
		return
	}

	// Add the courier to the store's couriers_ids array
	updateQuery := `
		UPDATE stores 
		SET couriers_ids = ARRAY_APPEND(couriers_ids, $1), updated_at = NOW()
		WHERE id = $2
	`
	_, err = utils.DB.Exec(updateQuery, courierUUID.String(), storeID)
	if err != nil {
		log.Println("Error adding courier to store:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Courier added successfully"})
}

// RemoveCourierFromStore godoc
// @Summary Remove a courier from a store
// @Description Remove a courier by user ID from a store
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param request body models.RemoveCourierRequest true "Courier User ID"
// @Success 200 {object} map[string]string "Courier removed successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/remove-courier/{id} [patch]
// @Tag Stores
func (sc *StoreController) RemoveCourierFromStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	// Parse and validate the request body
	var req models.RemoveCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CourierID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateQuery := `UPDATE stores SET couriers_ids = array_remove(couriers_ids, $1), updated_at = $2 WHERE id = $3`
	_, err := utils.DB.Exec(updateQuery, req.CourierID, time.Now(), storeID)
	if err != nil {
		log.Println("Error removing courier from store:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Courier removed successfully"})
}

// GetAllAdminsForStore godoc
// @Summary Get all admins for a store
// @Description Get detailed information about all admins for a given store
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {array} models.User "List of admin details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/get-admins/{id} [get]
// @Tag Stores
func (sc *StoreController) GetStoreAdminById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT admins_ids FROM stores WHERE id = $1`
	var adminIDs pq.StringArray

	err := utils.DB.QueryRow(query, storeID).Scan(&adminIDs)
	if err == sql.ErrNoRows {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error fetching admin IDs:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if len(adminIDs) == 0 {
		http.Error(w, "No admins found for this store", http.StatusNotFound)
		return
	}

	adminDetailsQuery := `
		SELECT u.id, u.name, u.email, u.phone, u.location, u.created_at
		FROM admins a
		JOIN users u ON a.user_id = u.id
		WHERE a.store_id = $1 AND a.id = ANY($2)
	`

	rows, err := utils.DB.Query(adminDetailsQuery, storeID, pq.Array(adminIDs))
	if err != nil {
		log.Println("Error fetching admin details:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var admins []models.User
	for rows.Next() {
		var admin models.User
		if err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Phone, &admin.Location, &admin.CreatedAt); err != nil {
			log.Println("Error scanning admin row:", err)
			continue
		}
		admins = append(admins, admin)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admins)
}

// GetAllCouriersForStore godoc
// @Summary Get all couriers for a store
// @Description Get detailed information about all couriers for a given store
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {array} models.Courier "List of courier details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Store not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/get-couriers/{id} [get]
// @Tag Stores
func (sc *StoreController) GetStoreCourierById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeID := vars["id"]

	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT couriers_ids FROM stores WHERE id = $1`
	var courierIDs pq.StringArray

	err := utils.DB.QueryRow(query, storeID).Scan(&courierIDs)
	if err == sql.ErrNoRows {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error fetching courier IDs:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if len(courierIDs) == 0 {
		http.Error(w, "No couriers found for this store", http.StatusNotFound)
		return
	}

	courierDetailsQuery := `
		SELECT u.id, u.name, u.email, u.phone, u.location, u.created_at,
		       c.vehicle_type, c.available, c.last_active_at, c.store_id
		FROM couriers c
		JOIN users u ON c.user_id = u.id
		WHERE c.store_id = $1 AND c.id = ANY($2)
	`

	rows, err := utils.DB.Query(courierDetailsQuery, storeID, pq.Array(courierIDs))
	if err != nil {
		log.Println("Error fetching courier details:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var couriers []models.Courier
	for rows.Next() {
		var courier models.Courier
		if err := rows.Scan(
			&courier.ID, &courier.Name, &courier.Email, &courier.Phone, &courier.Location, &courier.CreatedAt,
			&courier.VehicleType, &courier.Available, &courier.LastActiveAt, &courier.StoreId,
		); err != nil {
			log.Println("Error scanning courier row:", err)
			continue
		}
		couriers = append(couriers, courier)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(couriers)
}

// GetAllAdmins godoc
// @Summary Get all admins for a specific store
// @Description Get detailed information about all admins associated with a given store
// @Accept json
// @Produce json
// @Param storeId query string true "Store ID"
// @Success 200 {array} models.User "List of admin details"
// @Failure 400 {object} map[string]string "Invalid store ID"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/get-admins [get]
// @Tag Admins
func (sc *StoreController) GetAllAdminsForStore(w http.ResponseWriter, r *http.Request) {
	storeId := r.URL.Query().Get("storeId")
	if storeId == "" {
		http.Error(w, "Missing storeId", http.StatusBadRequest)
		return
	}

	adminDetailsQuery := `
		SELECT u.id, u.name, u.email, u.phone, u.location, u.created_at
		FROM admins a
		JOIN users u ON a.user_id = u.id
		WHERE a.store_id = $1
	`

	rows, err := utils.DB.Query(adminDetailsQuery, storeId)
	if err != nil {
		log.Println("Error fetching admin details:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var admins []models.User
	for rows.Next() {
		var admin models.User
		if err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Phone, &admin.Location, &admin.CreatedAt); err != nil {
			log.Println("Error scanning admin row:", err)
			continue
		}
		admins = append(admins, admin)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admins)
}

// GetAllCouriers godoc
// @Summary Get all couriers for a specific store
// @Description Get detailed information about all couriers associated with a given store
// @Accept json
// @Produce json
// @Param storeId query string true "Store ID"
// @Success 200 {array} models.Courier "List of courier details"
// @Failure 400 {object} map[string]string "Invalid store ID"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/get-couriers [get]
// @Tag Couriers
func (sc *StoreController) GetAllCouriersForStore(w http.ResponseWriter, r *http.Request) {
	storeId := r.URL.Query().Get("storeId")
	if storeId == "" {
		http.Error(w, "Missing storeId", http.StatusBadRequest)
		return
	}

	courierDetailsQuery := `
		SELECT u.id, u.name, u.email, u.phone, u.location, u.created_at,
		       c.vehicle_type, c.available, c.last_active_at, c.store_id
		FROM couriers c
		JOIN users u ON c.user_id = u.id
		WHERE c.store_id = $1
	`

	rows, err := utils.DB.Query(courierDetailsQuery, storeId)
	if err != nil {
		log.Println("Error fetching courier details:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var couriers []models.Courier
	for rows.Next() {
		var courier models.Courier
		if err := rows.Scan(
			&courier.ID, &courier.Name, &courier.Email, &courier.Phone, &courier.Location, &courier.CreatedAt,
			&courier.VehicleType, &courier.Available, &courier.LastActiveAt, &courier.StoreId,
		); err != nil {
			log.Println("Error scanning courier row:", err)
			continue
		}
		couriers = append(couriers, courier)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(couriers)
}

// GetAvailableCouriersByStore godoc
// @Summary Get available couriers by store ID
// @Description Retrieve a list of available couriers for a specific store.
// @Accept json
// @Produce json
// @Param store_id query string true "Store ID"
// @Success 200 {array} map[string]interface{} "List of available couriers"
// @Failure 400 {object} map[string]string "Invalid store ID"
// @Failure 404 {object} map[string]string "Store not found or no available couriers"
// @Failure 500 {object} map[string]string "Server error"
// @Router /stores/couriers/available [get]
// @Tag Store
func (ac *StoreController) GetAvailableCouriersByStore(w http.ResponseWriter, r *http.Request) {
	// Get the store ID from the query parameters
	storeID := r.URL.Query().Get("store_id")
	if storeID == "" {
		http.Error(w, "Store ID is required", http.StatusBadRequest)
		return
	}

	// Check if the store exists
	var storeExists bool
	checkStoreQuery := "SELECT EXISTS (SELECT 1 FROM stores WHERE id = $1)"
	err := utils.DB.QueryRow(checkStoreQuery, storeID).Scan(&storeExists)
	if err != nil {
		log.Println("Error checking store existence:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if !storeExists {
		http.Error(w, "Store not found", http.StatusNotFound)
		return
	}

	// Query to get all available couriers for the given store ID
	query := `
        SELECT 
            u.id, u.name, u.email, u.phone, u.location, u.created_at,
            c.vehicle_type, c.available, c.last_active_at, c.store_id
        FROM users u
        JOIN couriers c ON u.id = c.user_id
        WHERE c.store_id = $1 AND c.available = true
    `

	// Execute the query
	rows, err := utils.DB.Query(query, storeID)
	if err != nil {
		log.Println("Error retrieving available couriers:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Prepare a list to store available couriers
	var availableCouriers []map[string]interface{}

	// Iterate through the result set and add couriers to the list
	for rows.Next() {
		var user models.User
		var courier models.Courier

		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Phone, &user.Location, &user.CreatedAt,
			&courier.VehicleType, &courier.Available, &courier.LastActiveAt, &courier.StoreId,
		)
		if err != nil {
			log.Println("Error scanning courier data:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Append courier data to the list
		courierData := map[string]interface{}{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"phone":       user.Phone,
			"location":    user.Location,
			"created_at":  user.CreatedAt,
			"vehicleType": courier.VehicleType,
			"available":   courier.Available,
			"lastActive":  courier.LastActiveAt,
			"store_id":    courier.StoreId,
		}
		availableCouriers = append(availableCouriers, courierData)
	}

	// Check if there are no available couriers
	if len(availableCouriers) == 0 {
		http.Error(w, "No available couriers found", http.StatusNotFound)
		return
	}

	// Send the response with the list of available couriers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableCouriers)
}
