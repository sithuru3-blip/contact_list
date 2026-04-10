package models

// User model represents users table in DB
type User struct {
    ID       uint   `gorm:"primaryKey"`        // Primary Key (auto increment)
    Name     string `gorm:"not null"`          // User name (required)
    Email    string `gorm:"unique;not null"`   // Unique email (no duplicates)
    Password string `gorm:"not null"`          // Hashed password
}
