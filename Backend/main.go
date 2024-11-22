package main

import (
	"PTS/APIs"
	"PTS/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize logging and server details
	fmt.Println("Starting the server...")

	// Connect to the database
	utils.ConnectDB()

	// Initialize the router
	router := mux.NewRouter()

	// Wrap the router with CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"}, // Allow frontend to access
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler := corsHandler.Handler(router)

	// Register API routes
	APIs.RegisterAuthRoutes(router)
	APIs.RegisterItemRoutes(router)
	APIs.RegisterOrderRoutes(router)
	APIs.RegisterStoreRoutes(router)

	// Serve Swagger JSON
	router.Path("/swagger/doc.json").HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	}))

	// Serve Swagger UI
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

	// Automatically open the Swagger UI page in the default browser (Windows-specific)
	go func() {
		exec.Command("cmd", "/c", "start", "http://localhost:8080/swagger").Run()
	}()

	// Set up graceful shutdown handling
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		// Wait for interrupt signal to gracefully shut down the server
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		log.Println("Shutting down the server...")

		// Attempt to shut down gracefully
		if err := server.Shutdown(nil); err != nil {
			log.Fatal("Server Shutdown Failed:", err)
		}
		log.Println("Server gracefully stopped")
	}()

	// Start the server on port 8080
	log.Println("Server running on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error starting server: ", err)
	}
}
