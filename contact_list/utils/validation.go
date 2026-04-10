package utils

import (
	"regexp"
	"errors"
	"strings"
)

// PasswordValidationResult contains validation details
type PasswordValidationResult struct {
	IsValid bool
	Errors  []string   // validation result return
}

// ValidatePassword checks if password meets requirements:
// - Minimum 8 characters
// - At least one uppercase letter
// - At least one number
// - At least one special character
func ValidatePassword(password string) error {
	var validationErrors []string
	
	// Check minimum length
	if len(password) < 8 {
		validationErrors = append(validationErrors, "password must be at least 8 characters long")
	}
	
	// Check for uppercase letter
	uppercasePattern := `[A-Z]`
	if match, _ := regexp.MatchString(uppercasePattern, password); !match {
		validationErrors = append(validationErrors, "password must contain at least one uppercase letter")
	}
	
	// Check for number
	numberPattern := `[0-9]`
	if match, _ := regexp.MatchString(numberPattern, password); !match {
		validationErrors = append(validationErrors, "password must contain at least one number")
	}
	
	// Check for special character
	specialCharPattern := `[!@#~$%^&*()_+{}":;'?/>.<,\\|[\]\\=-]`
	if match, _ := regexp.MatchString(specialCharPattern, password); !match {
		validationErrors = append(validationErrors, "password must contain at least one special character")
	}
	
	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}
	
	return nil
}

// ValidatePasswordStrength returns detailed validation result
func ValidatePasswordStrength(password string) PasswordValidationResult {
	result := PasswordValidationResult{
		IsValid: true,
		Errors:  []string{},
	}
	
	if len(password) < 8 {
		result.IsValid = false
		result.Errors = append(result.Errors, "At least 8 characters")
	}
	
	uppercasePattern := `[A-Z]`
	if match, _ := regexp.MatchString(uppercasePattern, password); !match {
		result.IsValid = false
		result.Errors = append(result.Errors, "At least one uppercase letter")
	}
	
	numberPattern := `[0-9]`
	if match, _ := regexp.MatchString(numberPattern, password); !match {
		result.IsValid = false
		result.Errors = append(result.Errors, "At least one number")
	}
	
	specialCharPattern := `[!@#~$%^&*()_+{}":;'?/>.<,\\|[\]\\=-]`
	if match, _ := regexp.MatchString(specialCharPattern, password); !match {
		result.IsValid = false
		result.Errors = append(result.Errors, "At least one special character")
	}
	
	return result
}