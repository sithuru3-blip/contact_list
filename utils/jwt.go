package utils

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

// Secret key for signing JWT (keep secure in real apps)
var jwtKey = []byte("your_secret_key")

// Claims struct → data stored inside JWT
type Claims struct {
    UserID int `json:"user_id"`     // Custom field (user ID)
    jwt.StandardClaims              // Includes expiry, issuer, etc.
}

// GenerateToken creates JWT for a given user
func GenerateToken(userID int) (string, error) {

    // Set expiration time (15 minutes)
    expirationTime := time.Now().Add(15 * time.Minute)

    // Create claims object
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(), // Expiry time
        },
    }

    // Create token with HS256 signing method
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Return signed token string
    return token.SignedString(jwtKey)
}

// ValidateToken verifies token and extracts data
func ValidateToken(tokenString string) (*Claims, error) {

    claims := &Claims{}

    // Parse token with claims
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    // If invalid → return error
    if err != nil || !token.Valid {
        return nil, err
    }

    // Return decoded claims
    return claims, nil
}

// GenerateRefreshToken creates long-lived refresh token
func GenerateRefreshToken(userID int) (string, error) {

    expirationTime := time.Now().Add(7 * 24 * time.Hour)

    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString(jwtKey)
}