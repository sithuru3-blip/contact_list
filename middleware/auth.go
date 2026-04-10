package middleware

import (
    "contactlist/utils"
    "context"
    "net/http"
    "strings"
)

// Auth middleware checks JWT before accessing protected routes
func Auth(next http.HandlerFunc) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        // Get Authorization header - token
        authHeader := r.Header.Get("Authorization")

        // If missing → unauthorized
        if authHeader == "" {
            http.Error(w, "Missing token", 401)
            return
        }

        // Remove "Bearer " prefix
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Validate token using utils
        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", 401)
            return
        }

        // Store userID in request context
        ctx := context.WithValue(r.Context(), "userID", claims.UserID)

        // Continue to next handler
        next(w, r.WithContext(ctx))
    }
}
