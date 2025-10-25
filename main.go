package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StructA struct {
	FieldA string `form:"field_a"`
}

// Define a struct to map query parameters
type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}
type Student struct {
	Name string `uri:"name" binding:"required"`
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println("status = ", status)
	}
}

func main() {
	router := gin.Default()
	// router.Use(Logger())

	router.GET("/testing", startPage)
	router.GET("/:name", student)
	router.GET("/user", users)
	router.GET("/long_async", long_async)

	router.GET("/long_sync", long_sync)

	router.Run(":8080")
}
func users(c *gin.Context) {
	var person Person

	// Bind query parameters into the struct
	if err := c.ShouldBindQuery(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond with the extracted values
	c.JSON(http.StatusOK, gin.H{
		"name":    person.Name,
		"address": person.Address,
	})
}

func startPage(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	err := c.ShouldBind(&person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(person.Name)
	log.Println(person.Address)
	log.Println(person.Birthday)

	c.String(200, "Success")
}

func student(c *gin.Context) {

	var student Student
	if err := c.ShouldBindUri(&student); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": student.Name})

}

func long_async(c *gin.Context) {
	// create copy to be used inside the goroutine
	cCp := c.Copy()
	go func() {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// note that you are using the copied context "cCp", IMPORTANT
		log.Println("Done! in path " + cCp.Request.URL.Path)
	}()
	log.Println("in async")
}

func long_sync(c *gin.Context) {
	// simulate a long task with time.Sleep(). 5 seconds
	time.Sleep(5 * time.Second)

	// since we are NOT using a goroutine, we do not have to copy the context
	log.Println("Done! in path " + c.Request.URL.Path)
	log.Println("in async")

}
