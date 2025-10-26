package dto

import (
	"fmt"
	"strings"
	"time"
)

// CategoryRequestDTO is used to create or update a category.
type CategoryRequestDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=255"`
}

// Validate performs additional business logic validation
func (c *CategoryRequestDTO) Validate() error {
	// Trim spaces and check name
	c.Name = strings.TrimSpace(c.Name)
	if len(c.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(c.Name) > 50 {
		return fmt.Errorf("name must not exceed 50 characters")
	}

	// Check description length
	if len(c.Description) > 255 {
		return fmt.Errorf("description must not exceed 255 characters")
	}

	return nil
}

// CategoryResponseDTO represents a category returned in API responses.
type CategoryResponseDTO struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
