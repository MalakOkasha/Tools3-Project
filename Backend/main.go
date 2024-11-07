package main

import (
	UserAPIs "PTS/APIs"
	"PTS/utils"
	"fmt"
	"log"
	"net/http"
	"os/exec"

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
// @BasePath        /
func main() {
	fmt.Println("Starting the server...")

	// Connect to the database
	utils.ConnectDB()

	// Initialize the router
	router := mux.NewRouter()

	// Wrap the router with CORS middleware
	handler := cors.Default().Handler(router)

	// Register API routes
	UserAPIs.RegisterAuthRoutes(router)

	// Serve Swagger JSON
	router.Path("/swagger/doc.json").HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	}))

	// Serve the Swagger UI HTML page
	router.Path("/swagger").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlContent := `<!DOCTYPE html>
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

	// Automatically open the Swagger UI page in the default browser
	go func() {
		exec.Command("cmd", "/c", "start", "http://localhost:8080/swagger").Run()
	}()

	// Start the server on port 8080
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
