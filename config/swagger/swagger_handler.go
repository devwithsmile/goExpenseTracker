package swagger

import (
	"encoding/json"
	"net/http"
	"os"

	"goExpenseTracker/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

// CustomSwaggerHandler modifies the swagger spec to include multiple servers
func CustomSwaggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the swagger spec
		swagger, err := swag.ReadDoc(docs.SwaggerInfo.InstanceName())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read swagger spec"})
			return
		}

		// Parse the swagger spec
		var spec map[string]interface{}
		if err := json.Unmarshal([]byte(swagger), &spec); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse swagger spec"})
			return
		}

		// Get server URLs from environment variables with defaults
		localURL := os.Getenv("SWAGGER_LOCAL_URL")
		if localURL == "" {
			localURL = "http://localhost:8080/api"
		}

		prodURL := os.Getenv("SWAGGER_PROD_URL")
		if prodURL == "" {
			prodURL = "https://goexpensetracker.onrender.com/api"
		}

		// Add servers array for OpenAPI-style server selection (Swagger UI 3+ supports this)
		spec["servers"] = []map[string]interface{}{
			{
				"url":         localURL,
				"description": "Local Development Server",
			},
			{
				"url":         prodURL,
				"description": "Production Server (Render)",
			},
		}

		// Marshal with proper formatting
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSON(http.StatusOK, spec)
	}
}
