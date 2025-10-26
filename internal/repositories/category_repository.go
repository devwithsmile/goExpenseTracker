package repositories

import (
	"goExpenseTracker/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetAll(offset int, limit int, nameFilter string) ([]models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetAll fetches categories with pagination and optional filters
func (r *categoryRepository) GetAll(offset int, limit int, nameFilter string) ([]models.Category, error) {
	var categories []models.Category

	query := r.db.Model(&models.Category{})

	if nameFilter != "" {
		query = query.Where("name ILIKE ?", "%"+nameFilter+"%")
	}

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
