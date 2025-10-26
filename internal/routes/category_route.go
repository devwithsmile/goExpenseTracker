package routes

import (
	"goExpenseTracker/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.RouterGroup, categoryHandler *handlers.CategoryHandler) {
	v1 := router.Group("/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetAllCategorys)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}
	}
}
