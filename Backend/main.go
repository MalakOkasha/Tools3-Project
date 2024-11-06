package main

import (
	UserAPIs "PTS/APIs"
	//"PTS/controllers"
	//"PTS/models"
	"PTS/utils"
	//"bytes"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"net/http/httptest"
	"os/exec"

	_ "PTS/docs" // Import the generated Swagger documentation

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// @title           Package Tracking System (PTS-OpenShift) phase 0
// @version         1.0
// @description     This is a sample API for user registration and login.
// @termsOfService  http://example.com/terms/
// @contact.name    Abdulrahman Hijazy
// @contact.url     https://www.linkedin.com/in/abdulrahmanhijazy
// @contact.email   abdulrahman.hijazy.a@gmail.com
// @license.name    Cairo University
// @license.url     Project Repo link
// @host            localhost:8080
// @BasePath       /
func main() {
	fmt.Println("test 1")
	utils.ConnectDB()

	router := mux.NewRouter()
	fmt.Println("test 2")

	// Wrap your router with CORS middleware
	handler := cors.Default().Handler(router)

	UserAPIs.RegisterAuthRoutes(router)

	// Serve Swagger JSON
	router.Path("/swagger/doc.json").HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	}))

	// Serve the Swagger UI HTML page
	router.Path("/swagger").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Swagger UI</title>
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/swagger-ui.css">
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/swagger-ui-bundle.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/swagger-ui-standalone-preset.js"></script>
		</head>
		<body>
			<div id="swagger-ui"></div>
			<script>
				const ui = SwaggerUIBundle({
					url: '/swagger/doc.json', // Swagger JSON file URL
					dom_id: '#swagger-ui',
					presets: [
						SwaggerUIBundle.presets.apis,
						SwaggerUIStandalonePreset
					],
				});
			</script>
		</body>
		</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})

	// to register Statically

	// // Create a static user object for registration
	// staticUser := models.User{
	// 	Name:     "test",
	// 	Email:    "tdsddtagfqsfssaadadt@gmail.com",
	// 	Phone:    "01010101010",
	// 	Password: "test",
	// }

	// authController := &controllers.AuthController{}
	// if err := registerStaticUser(authController, staticUser); err != nil {
	// 	log.Fatalf("Error registering static user: %v", err)
	// }

	// fmt.Println("Static user registered successfully")

	// Automatically open the Swagger UI page
	go func() {
		exec.Command("cmd", "/c", "start", "http://localhost:8080/swagger").Run()
	}()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

	// // Function to register a static user
	// func registerStaticUser(authController *controllers.AuthController, user models.User) error {
	// 	reqBody, err := json.Marshal(controllers.RegisterRequest{
	// 		Name:     user.Name,
	// 		Email:    user.Email,
	// 		Phone:    user.Phone,
	// 		Password: user.Password,
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to marshal request body: %w", err)
	// 	}

	// 	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	// 	req.Header.Set("Content-Type", "application/json")
	// 	rr := httptest.NewRecorder()
	// 	authController.Register(rr, req)

	// 	if rr.Code != http.StatusCreated {
	// 		return fmt.Errorf("failed to register static user: %s", rr.Body.String())
	// 	}

	// 	return nil
	// }

}
