package routes

import (
	"goExpenseTracker/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupExpenseRoutes(router *gin.RouterGroup, expenseHandler *handlers.ExpenseHandler) {
	v1 := router.Group("/v1")
	{
		expenses := v1.Group("/expenses")
		{
			expenses.POST("", expenseHandler.CreateExpense)
			expenses.GET("", expenseHandler.GetAllExpenses)
			expenses.GET("/:id", expenseHandler.GetExpenseByID)
			expenses.PUT("/:id", expenseHandler.UpdateExpense)
			expenses.DELETE("/:id", expenseHandler.DeleteExpense)
		}
	}
}
