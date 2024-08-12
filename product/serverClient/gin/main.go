package main

import (
	"fmt"
	routes "github.com/SaYaku64/hillel/product/serverClient/internal"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {

	routes.Homework()
	r := gin.Default()

	r.LoadHTMLGlob("front/templates/*")
	r.Static("/static", "./front/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"version": fmt.Sprintf("%v", time.Now().Unix()),
			"title":   "Title of site",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/alert", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, everyone!")
		})
		v1.POST("/post", postHandler)

		v1.GET("/getSessionID", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%v", time.Now().Unix()))
		})

	}

	// Запуск сервера
	r.Run(":9095")
}

func postHandler(c *gin.Context) {
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
}
