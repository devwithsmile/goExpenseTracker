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

type ExpenseHandler struct {
	ExpenseService services.ExpenseService
}

// NewExpenseHandler initializes a new ExpenseHandler
func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		ExpenseService: service,
	}
}

// formatValidationError converts validation errors into user-friendly messages
func formatValidationError(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := strings.ToLower(strings.Replace(e.Field(), "ID", "_id", 1))

			switch e.Tag() {
			case "required":
				return field + " is required"
			case "min":
				if e.Field() == "CategoryID" {
					return "category_id must be greater than 0"
				}
				return field + " must be at least " + e.Param()
			case "gt":
				return field + " must be greater than 0"
			case "max":
				return field + " must not exceed " + e.Param() + " characters"
			default:
				return field + " is invalid"
			}
		}
	}
	return err.Error()
}

// CreateExpense godoc
// @Summary      Create a new expense
// @Description  Add a new expense entry including amount, category, and description
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        expense  body      dto.ExpenseRequestDTO  true  "Expense Data"
// @Success      201       {object}  dto.ExpenseResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /v1/expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req dto.ExpenseRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	// Additional validation
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdExpense, err := h.ExpenseService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdExpense)
}

// GetAllExpenses godoc
// @Summary      Get all expenses
// @Description  Retrieve all recorded expenses with pagination and filtering
// @Tags         expenses
// @Param        description  query  string  false  "Filter by description (partial match)"
// @Param        category_id  query  int     false  "Filter by category ID"
// @Produce      json
// @Param        offset       query  int     false  "Offset for pagination" default(0)
// @Param        limit        query  int     false  "Limit for pagination" default(10)
// @Success      200  {array}   dto.ExpenseResponseDTO
// @Failure      500  {object}  map[string]string
// @Router       /v1/expenses [get]
func (h *ExpenseHandler) GetAllExpenses(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	descriptionFilter := c.DefaultQuery("description", "")
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))

	expenses, err := h.ExpenseService.GetAll(offset, limit, descriptionFilter, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// GetExpenseByID godoc
// @Summary      Get expense by ID
// @Description  Retrieve details of a specific expense by its ID
// @Tags         expenses
// @Produce      json
// @Param        id   path      int  true  "Expense ID"
// @Success      200  {object}  dto.ExpenseResponseDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /v1/expenses/{id} [get]
func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	expense, err := h.ExpenseService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateExpense godoc
// @Summary      Update expense
// @Description  Update details of a specific expense
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        id        path      int                    true  "Expense ID"
// @Param        expense   body      dto.ExpenseRequestDTO  true  "Updated Expense Data"
// @Success      200       {object}  dto.ExpenseResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /v1/expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	var req dto.ExpenseRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	// Additional validation
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedExpense, err := h.ExpenseService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedExpense)
}

// DeleteExpense godoc
// @Summary      Delete expense
// @Description  Remove a specific expense entry by ID
// @Tags         expenses
// @Param        id   path  int  true  "Expense ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	err = h.ExpenseService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
