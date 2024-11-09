package APIs

import (
	"PTS/controllers"

	"github.com/gorilla/mux"
)

// RegisterItemRoutes registers item-related routes
func RegisterItemRoutes(router *mux.Router) {
	itemController := &controllers.ItemController{}

	// Routes for Items
	router.HandleFunc("/items/add", itemController.AddItem).Methods("POST")
	router.HandleFunc("/items/get/{id}", itemController.GetItemByID).Methods("GET")
	router.HandleFunc("/items/list/{store_id}", itemController.ListItemsByStoreID).Methods("GET")
	router.HandleFunc("/items/delete/{id}", itemController.DeleteItem).Methods("DELETE")
}
