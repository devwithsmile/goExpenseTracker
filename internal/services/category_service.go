package services

import (
	dto "goExpenseTracker/internal/DTOs"
	"goExpenseTracker/internal/models"
	"goExpenseTracker/internal/repositories"
	"time"
)

type CategoryService interface {
	Create(req dto.CategoryRequestDTO) (dto.CategoryResponseDTO, error)
	GetAll(offset, limit int, nameFilter string) ([]dto.CategoryResponseDTO, error)
	GetByID(id int) (dto.CategoryResponseDTO, error)
	Update(id int, req dto.CategoryRequestDTO) (dto.CategoryResponseDTO, error)
	Delete(id int) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// Create category
func (s *categoryService) Create(req dto.CategoryRequestDTO) (dto.CategoryResponseDTO, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Persist model
	err := s.repo.Create(&category)
	if err != nil {
		return dto.CategoryResponseDTO{}, err
	}

	return s.toResponseDTO(category), nil
}

// Get all categories
func (s *categoryService) GetAll(offset, limit int, nameFilter string) ([]dto.CategoryResponseDTO, error) {
	categories, err := s.repo.GetAll(offset, limit, nameFilter)
	if err != nil {
		return []dto.CategoryResponseDTO{}, err
	}

	responses := make([]dto.CategoryResponseDTO, 0)
	for _, category := range categories {
		responses = append(responses, s.toResponseDTO(category))
	}
	return responses, nil
}

// Get single category
func (s *categoryService) GetByID(id int) (dto.CategoryResponseDTO, error) {
	categoryPtr, err := s.repo.GetByID(uint(id))
	if err != nil {
		return dto.CategoryResponseDTO{}, err
	}
	return s.toResponseDTO(*categoryPtr), nil
}

// Update category
func (s *categoryService) Update(id int, req dto.CategoryRequestDTO) (dto.CategoryResponseDTO, error) {
	existing, err := s.repo.GetByID(uint(id))
	if err != nil {
		return dto.CategoryResponseDTO{}, err
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.UpdatedAt = time.Now()

	err = s.repo.Update(existing)
	if err != nil {
		return dto.CategoryResponseDTO{}, err
	}

	return s.toResponseDTO(*existing), nil
}

// Delete category
func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(uint(id))
}

// Private helper for mapping model → DTO
func (s *categoryService) toResponseDTO(category models.Category) dto.CategoryResponseDTO {
	return dto.CategoryResponseDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}
