package main

import (
    "contactlist/db"
    "contactlist/handlers"
    "contactlist/middleware"
    "contactlist/models"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

// Add CORS middleware
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    // Connect to database
    db.Connect()

    // Auto create/update tables
    db.DB.AutoMigrate(&models.User{}, &models.Contact{})
    
    // Initialize router
    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
    r.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")
    r.HandleFunc("/refresh", handlers.RefreshToken).Methods("POST", "OPTIONS")
    r.HandleFunc("/logout", handlers.Logout).Methods("POST", "OPTIONS")

    // Protected routes (JWT required)
    r.HandleFunc("/contacts", middleware.Auth(handlers.AddContact)).Methods("POST", "OPTIONS")
    r.HandleFunc("/contacts", middleware.Auth(handlers.ListContacts)).Methods("GET", "OPTIONS")
    r.HandleFunc("/contacts/{id}", middleware.Auth(handlers.UpdateContact)).Methods("PUT", "OPTIONS")
    r.HandleFunc("/contacts/{id}", middleware.Auth(handlers.DeleteContact)).Methods("DELETE", "OPTIONS")
    r.HandleFunc("/contacts/search", middleware.Auth(handlers.SearchContacts)).Methods("GET", "OPTIONS")

    // Apply CORS middleware
    handler := enableCORS(r)

    log.Println("Server running on :8080")

    // Start server with CORS enabled
    log.Fatal(http.ListenAndServe(":8080", handler))
}