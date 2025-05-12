package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Server is healthy!",
		ResponseCode: http.StatusOK,
	} 

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	// Create a new router
	r := mux.NewRouter()
	r.Use(RateLimiter)

	// Define routes
	r.HandleFunc("/api/health", healthCheck).Methods("GET")
	r.HandleFunc("/api/application", ApplicationHandler).Methods("POST")
	r.HandleFunc("/api/reservation", ReservationHandler).Methods("POST")

	// Setup CORS
	reactUrl := os.Getenv("REACT_URL")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{reactUrl}, // React app's URL
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := c.Handler(r)

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
