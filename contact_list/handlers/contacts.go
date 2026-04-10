package handlers

import (
    "contactlist/db"
    "contactlist/models"
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
)

// AddContact creates a new contact
func AddContact(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

    var contact models.Contact
    if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
        http.Error(w, "Invalid request body", 400)
        return
    }

    contact.UserID = uint(userID)

    if err := db.DB.Create(&contact).Error; err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contact)
}

// ListContacts returns all contacts for logged-in user
func ListContacts(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

    var contacts []models.Contact
    db.DB.Where("user_id = ?", userID).Find(&contacts)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contacts)
}

// UpdateContact updates a specific contact - FIXED VERSION
func UpdateContact(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    id := mux.Vars(r)["id"]

    // First find the existing contact
    var existingContact models.Contact
    if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&existingContact).Error; err != nil {
        http.Error(w, "Contact not found", 404)
        return
    }

    // Decode updated data
    var updatedData models.Contact
    if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
        http.Error(w, "Invalid request body", 400)
        return
    }

    // Update only the fields that are provided
    if updatedData.Name != "" {
        existingContact.Name = updatedData.Name
    }
    if updatedData.Phone != "" {
        existingContact.Phone = updatedData.Phone
    }
    if updatedData.Email != "" {
        existingContact.Email = updatedData.Email
    }
    if updatedData.Address != "" {
        existingContact.Address = updatedData.Address
    }
    if updatedData.Tags != "" {
        existingContact.Tags = updatedData.Tags
    }
    if updatedData.Notes != "" {
        existingContact.Notes = updatedData.Notes
    }

    // Save the updated contact
    if err := db.DB.Save(&existingContact).Error; err != nil {
        http.Error(w, "Failed to update contact", 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Contact updated successfully",
        "contact": existingContact,
    })
}

// DeleteContact deletes a contact
func DeleteContact(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    id := mux.Vars(r)["id"]

    result := db.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Contact{})
    
    if result.RowsAffected == 0 {
        http.Error(w, "Contact not found", 404)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Contact deleted successfully"})
}

// SearchContacts finds contacts by name - FIXED VERSION
func SearchContacts(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)
    q := r.URL.Query().Get("q")

    var contacts []models.Contact

    if q == "" {
        // If no search query, return all contacts
        db.DB.Where("user_id = ?", userID).Find(&contacts)
    } else {
        // Case-insensitive search by name
        db.DB.Where("user_id = ? AND LOWER(name) LIKE LOWER(?)", userID, "%"+q+"%").Find(&contacts)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(contacts)
}