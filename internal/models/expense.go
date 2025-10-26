package models

import (
	"time"
)

type Expense struct {
	ID          int       `json:"id" db:"id"`
	CategoryID  int       `json:"category_id" db:"category_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
