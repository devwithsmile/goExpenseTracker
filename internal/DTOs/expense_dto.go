package dto

import (
	"fmt"
	"time"
)

// ExpenseRequestDTO is for creating or updating an expense.
type ExpenseRequestDTO struct {
	CategoryID  int     `json:"category_id" binding:"min=1"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"max=255"`
	Date        string  `json:"date" binding:"required" example:"12-12-2025" format:"date"` // Format: dd-mm-yyyy or yyyy-mm-dd
}

// Validate performs additional business logic validation
func (e *ExpenseRequestDTO) Validate() error {
	// Check if date is provided
	if e.Date == "" {
		return fmt.Errorf("date is required")
	}

	// Parse the date string
	var parsedDate time.Time
	var err error

	// Try parsing dd-mm-yyyy format
	parsedDate, err = time.Parse("02-01-2006", e.Date)
	if err != nil {
		// Try ISO format as fallback
		parsedDate, err = time.Parse("2006-01-02", e.Date)
		if err != nil {
			return fmt.Errorf("invalid date format, expected dd-mm-yyyy or yyyy-mm-dd")
		}
	}

	// Check if amount is positive (already in binding, but double-check)
	if e.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// Check if date is not in the past
	today := time.Now().Truncate(24 * time.Hour)
	expenseDate := parsedDate.Truncate(24 * time.Hour)

	if expenseDate.Before(today) {
		return fmt.Errorf("date cannot be in the past")
	}

	return nil
}

// ParseDate parses the date string into time.Time
func (e *ExpenseRequestDTO) ParseDate() (time.Time, error) {
	// Try parsing dd-mm-yyyy format
	t, err := time.Parse("02-01-2006", e.Date)
	if err != nil {
		// Try ISO format as fallback
		t, err = time.Parse("2006-01-02", e.Date)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid date format, expected dd-mm-yyyy or yyyy-mm-dd")
		}
	}
	return t, nil
}

// ExpenseResponseDTO represents the expense data sent to the client.
type ExpenseResponseDTO struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name,omitempty"` // Category name resolved from relationship
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	Date         string  `json:"date"` // Format: yyyy-mm-dd
}
