package APIs

import (
	"PTS/controllers"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(router *mux.Router) {
	userController := &controllers.UserController{}
	courierController := &controllers.CourierController{}
	adminController := &controllers.AdminController{}
	ownerController := &controllers.OwnerController{}

	// Routes for Normal Users
	router.HandleFunc("/users/register", userController.Register).Methods("POST")
	router.HandleFunc("/users/login", userController.Login).Methods("POST")

	// Routes for Courier Users
	router.HandleFunc("/couriers/register", courierController.CourierRegister).Methods("POST")
	router.HandleFunc("/couriers/login", courierController.CourierLogin).Methods("POST")

	// Routes for Admin Users
	router.HandleFunc("/admins/register", adminController.AdminRegister).Methods("POST")
	router.HandleFunc("/admins/login", adminController.AdminLogin).Methods("POST")

	// Routes for Owner Users
	router.HandleFunc("/owners/register", ownerController.OwnerRegister).Methods("POST")
	router.HandleFunc("/owners/login", ownerController.OwnerLogin).Methods("POST")
}
