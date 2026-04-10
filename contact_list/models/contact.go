package models

// Contact model represents contacts table
type Contact struct {
    ID      uint   `gorm:"primaryKey"` // Primary Key
    UserID  uint   `gorm:"index"`      // Foreign key (belongs to User)
    Name    string                     // Contact name
    Phone   string                     // Phone number
    Email   string                     // Email
    Address string                     // Address
    Tags    string                     // Tags (friend, work, etc.)
    Notes   string                     // Extra notes
}
