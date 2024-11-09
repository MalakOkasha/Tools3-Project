package APIs

import (
	"PTS/controllers"

	"github.com/gorilla/mux"
)

// RegisterOrderRoutes registers order-related routes
func RegisterOrderRoutes(router *mux.Router) {
	orderController := &controllers.OrderController{}

	// Routes for Orders
	router.HandleFunc("/orders/add", orderController.AddOrder).Methods("POST")
	router.HandleFunc("/orders/{order_id}", orderController.GetOrderDetails).Methods("GET")
	router.HandleFunc("/orders/list/user/{user_id}", orderController.ListOrdersByUserID).Methods("GET")
	router.HandleFunc("/orders/list/store/{store_id}", orderController.ListOrdersByStoreID).Methods("GET")
	router.HandleFunc("/orders/list/courier/{courier_id}", orderController.ListOrdersByCourierID).Methods("GET")
	router.HandleFunc("/orders/update/{id}", orderController.UpdateOrderStatus).Methods("PATCH")
	router.HandleFunc("/orders/delete/{id}", orderController.DeleteOrder).Methods("DELETE")
	router.HandleFunc("/orders/cancel/{id}", orderController.CancelOrder).Methods("PATCH")
}
