package handlers

import (
    "contactlist/db"
    "contactlist/models"
    "contactlist/utils"
    "encoding/json"
    "net/http"

    "golang.org/x/crypto/bcrypt"
)

// Register with enhanced password validation
func Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    
    // Decode request body
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", 400)
        return
    }
    
    // Validate email
    if user.Email == "" {
        http.Error(w, "Email is required", 400)
        return
    }
    
    // Validate name
    if user.Name == "" {
        http.Error(w, "Name is required", 400)
        return
    }
    
    // Validate password with comprehensive checks
    if err := utils.ValidatePassword(user.Password); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    
    // Check if user already exists
    var existingUser models.User
    if err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        http.Error(w, "Email already registered", 400)
        return
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error processing password", 500)
        return
    }
    user.Password = string(hashedPassword)
    
    // Create user
    if err := db.DB.Create(&user).Error; err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    
    // Don't return password in response
    user.Password = ""
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "User registered successfully",
        "user": map[string]interface{}{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}

// Login (unchanged)
func Login(w http.ResponseWriter, r *http.Request) {
    var input models.User
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", 400)
        return
    }
    
    var user models.User
    
    if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        http.Error(w, "Invalid email or password", 401)
        return
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        http.Error(w, "Invalid email or password", 401)
        return
    }
    
    // Generate tokens
    accessToken, err := utils.GenerateToken(int(user.ID))
    if err != nil {
        http.Error(w, "Error generating access token", 500)
        return
    }
    
    refreshToken, err := utils.GenerateRefreshToken(int(user.ID))
    if err != nil {
        http.Error(w, "Error generating refresh token", 500)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}

// RefreshToken (unchanged)
func RefreshToken(w http.ResponseWriter, r *http.Request) {
    var body struct {
        RefreshToken string `json:"refresh_token"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "Invalid request body", 400)
        return
    }
    
    claims, err := utils.ValidateToken(body.RefreshToken)
    if err != nil {
        http.Error(w, "Invalid refresh token", 401)
        return
    }
    
    newAccessToken, err := utils.GenerateToken(claims.UserID)
    if err != nil {
        http.Error(w, "Token error", 500)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newAccessToken,
    })
}

// Logout (unchanged)
func Logout(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Logout success",
    })
}