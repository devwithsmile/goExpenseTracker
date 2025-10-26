package repositories

import (
	"goExpenseTracker/internal/models"

	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	GetAll(offset int, limit int, descriptionFilter string, categoryID int) ([]models.Expense, error)
	GetByID(id uint) (*models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uint) error
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

// GetAll fetches expenses with pagination and optional filters
func (r *expenseRepository) GetAll(offset int, limit int, descriptionFilter string, categoryID int) ([]models.Expense, error) {
	var expenses []models.Expense

	query := r.db.Model(&models.Expense{})

	if descriptionFilter != "" {
		query = query.Where("description ILIKE ?", "%"+descriptionFilter+"%")
	}

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) GetByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}
