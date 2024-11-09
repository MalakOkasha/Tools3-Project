package APIs

import (
	"PTS/controllers"

	"github.com/gorilla/mux"
)

// RegisterStoreRoutes registers store-related routes
func RegisterStoreRoutes(router *mux.Router) {
	storeController := &controllers.StoreController{}

	// Routes for Stores
	router.HandleFunc("/stores/add", storeController.AddStore).Methods("POST")
	router.HandleFunc("/stores/get/{id}", storeController.GetStoreByID).Methods("GET")
	router.HandleFunc("/stores/list", storeController.ListStores).Methods("GET")
	router.HandleFunc("/stores/delete/{id}", storeController.DeleteStore).Methods("DELETE")
	router.HandleFunc("/stores/add-admin/{id}", storeController.AddAdminToStore).Methods("PATCH")
	router.HandleFunc("/stores/remove-admin/{id}", storeController.RemoveAdminFromStore).Methods("PATCH")
	router.HandleFunc("/stores/add-courier/{id}", storeController.AddCourierToStore).Methods("PATCH")
	router.HandleFunc("/stores/remove-courier/{id}", storeController.RemoveCourierFromStore).Methods("PATCH")
	router.HandleFunc("/stores/get-admins/{id}", storeController.GetStoreAdminById).Methods("GET")
	router.HandleFunc("/stores/get-couriers/{id}", storeController.GetStoreCourierById).Methods("GET")
	router.HandleFunc("/stores/get-admins", storeController.GetAllAdminsForStore).Methods("GET")
	router.HandleFunc("/stores/get-couriers", storeController.GetAllCouriersForStore).Methods("GET")
	router.HandleFunc("/stores/couriers/available", storeController.GetAvailableCouriersByStore).Methods("GET")

}
