package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func simple() {
	r := gin.Default()

	// Маршрут для GET-запиту
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, everyone!")
	})

	// Маршрут для POST-запиту з JSON-параметрами
	r.POST("/post", func(c *gin.Context) {
		var json struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Received",
			"name":    json.Name,
			"age":     json.Age,
		})
	})

	// Маршрут з параметрами URL
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello, %s!", name)
	})

	// Маршрут з параметрами запиту
	r.GET("/query", func(c *gin.Context) {
		name := c.DefaultQuery("name", "Guest")
		c.String(http.StatusOK, "Hello, %s!", name)
	})

	// Маршрут з обробкою форм
	r.POST("/form", formHandler)

	// Запуск сервера
	r.Run(":9090")
}

func formHandler(c *gin.Context) {
	name := c.PostForm("name")
	age := c.PostForm("age")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"age":  age,
	})
}
