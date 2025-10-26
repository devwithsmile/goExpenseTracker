package Logger

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// bodyLogWriter helps capture and store the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Wrap the response writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// After request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		responseBody := blw.body.String()

		log.Printf("\n---- Request Log ----\n")
		log.Printf("Client IP: %s", clientIP)
		log.Printf("Path: %s | Method: %s", c.Request.URL.Path, c.Request.Method)
		log.Printf("Request Body: %s", string(reqBody))
		log.Printf("Status Code: %d", statusCode)
		log.Printf("Latency: %v", latency)
		log.Printf("Response: %s\n", responseBody)
		log.Println("----------------------")
	}
}
