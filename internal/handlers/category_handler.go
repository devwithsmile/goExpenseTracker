package handlers

import (
	"net/http"
	"strconv"
	"strings"

	dto "goExpenseTracker/internal/DTOs"
	"goExpenseTracker/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	CategoryService services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: service,
	}
}

// formatCategoryValidationError converts validation errors into user-friendly messages
func formatCategoryValidationError(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := strings.ToLower(e.Field())

			switch e.Tag() {
			case "required":
				return field + " is required"
			case "min":
				return field + " must be at least " + e.Param() + " characters"
			case "max":
				return field + " must not exceed " + e.Param() + " characters"
			default:
				return field + " is invalid"
			}
		}
	}
	return err.Error()
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a category by providing name and description
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      dto.CategoryRequestDTO  true  "Category Data"
// @Success      201       {object}  dto.CategoryResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CategoryRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatCategoryValidationError(err)})
		return
	}

	createdCategory, err := h.CategoryService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

// GetAllCategorys godoc
// @Summary      Get all categories
// @Description  Retrieve list of all categories with pagination and filtering
// @Tags         categories
// @Produce      json
// @Param        name    query  string  false  "Filter by category name (partial match)"
// @Param        offset  query  int     false  "Offset for pagination" default(0)
// @Param        limit   query  int     false  "Limit for pagination" default(10)
// @Success      200  {array}   dto.CategoryResponseDTO
// @Failure      500  {object}  map[string]string
// @Router       /v1/categories [get]
func (h *CategoryHandler) GetAllCategorys(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.DefaultQuery("name", "")

	categories, err := h.CategoryService.GetAll(offset, limit, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Retrieve category details by its ID
// @Tags         categories
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  dto.CategoryResponseDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.CategoryService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Modify category name or description by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id        path      int                   true  "Category ID"
// @Param        category  body      dto.CategoryRequestDTO true  "Updated Category Data"
// @Success      200       {object}  dto.CategoryResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req dto.CategoryRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatCategoryValidationError(err)})
		return
	}

	updatedCategory, err := h.CategoryService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Remove category by its ID
// @Tags         categories
// @Param        id   path  int  true  "Category ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = h.CategoryService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
