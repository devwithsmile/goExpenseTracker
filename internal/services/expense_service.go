package services

import (
	"fmt"
	"time"

	dto "goExpenseTracker/internal/DTOs"
	"goExpenseTracker/internal/models"
	"goExpenseTracker/internal/repositories"
)

type ExpenseService interface {
	Create(req dto.ExpenseRequestDTO) (dto.ExpenseResponseDTO, error)
	GetAll(offset int, limit int, descriptionFilter string, categoryID int) ([]dto.ExpenseResponseDTO, error)
	GetByID(id int) (dto.ExpenseResponseDTO, error)
	Update(id int, req dto.ExpenseRequestDTO) (dto.ExpenseResponseDTO, error)
	Delete(id int) error
}

type expenseService struct {
	expenseRepo  repositories.ExpenseRepository
	categoryRepo repositories.CategoryRepository
}

func NewExpenseService(expenseRepo repositories.ExpenseRepository, categoryRepo repositories.CategoryRepository) ExpenseService {
	return &expenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

// Create a new expense
func (s *expenseService) Create(req dto.ExpenseRequestDTO) (dto.ExpenseResponseDTO, error) {
	// Verify category exists
	_, err := s.categoryRepo.GetByID(uint(req.CategoryID))
	if err != nil {
		return dto.ExpenseResponseDTO{}, fmt.Errorf("category not found")
	}

	// Parse the date
	parsedDate, err := req.ParseDate()
	if err != nil {
		return dto.ExpenseResponseDTO{}, err
	}

	expense := models.Expense{
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        parsedDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.expenseRepo.Create(&expense)
	if err != nil {
		return dto.ExpenseResponseDTO{}, err
	}

	return s.toResponseDTO(expense), nil
}

// Get all expenses
func (s *expenseService) GetAll(offset int, limit int, descriptionFilter string, categoryID int) ([]dto.ExpenseResponseDTO, error) {
	expenses, err := s.expenseRepo.GetAll(offset, limit, descriptionFilter, categoryID)
	if err != nil {
		return []dto.ExpenseResponseDTO{}, err
	}

	response := make([]dto.ExpenseResponseDTO, 0)
	for _, expense := range expenses {
		response = append(response, s.toResponseDTO(expense))
	}

	return response, nil
}

// Get expense by ID
func (s *expenseService) GetByID(id int) (dto.ExpenseResponseDTO, error) {
	expensePtr, err := s.expenseRepo.GetByID(uint(id))
	if err != nil {
		return dto.ExpenseResponseDTO{}, fmt.Errorf("expense not found")
	}
	return s.toResponseDTO(*expensePtr), nil
}

// Update existing expense
func (s *expenseService) Update(id int, req dto.ExpenseRequestDTO) (dto.ExpenseResponseDTO, error) {
	expense, err := s.expenseRepo.GetByID(uint(id))
	if err != nil {
		return dto.ExpenseResponseDTO{}, fmt.Errorf("expense not found")
	}

	// Verify category exists if changed
	if expense.CategoryID != req.CategoryID {
		_, err := s.categoryRepo.GetByID(uint(req.CategoryID))
		if err != nil {
			return dto.ExpenseResponseDTO{}, fmt.Errorf("category not found")
		}
	}

	// Parse the date
	parsedDate, err := req.ParseDate()
	if err != nil {
		return dto.ExpenseResponseDTO{}, err
	}

	expense.Amount = req.Amount
	expense.CategoryID = req.CategoryID
	expense.Description = req.Description
	expense.Date = parsedDate
	expense.UpdatedAt = time.Now()

	err = s.expenseRepo.Update(expense)
	if err != nil {
		return dto.ExpenseResponseDTO{}, err
	}

	return s.toResponseDTO(*expense), nil
}

// Delete expense by ID
func (s *expenseService) Delete(id int) error {
	return s.expenseRepo.Delete(uint(id))
}

// Helper: Convert model → Response DTO
func (s *expenseService) toResponseDTO(expense models.Expense) dto.ExpenseResponseDTO {
	// Fetch category name
	var categoryName string
	if category, err := s.categoryRepo.GetByID(uint(expense.CategoryID)); err == nil {
		categoryName = category.Name
	}

	return dto.ExpenseResponseDTO{
		ID:           expense.ID,
		Amount:       expense.Amount,
		CategoryID:   expense.CategoryID,
		CategoryName: categoryName,
		Description:  expense.Description,
		Date:         expense.Date.Format("2006-01-02"),
	}
}
