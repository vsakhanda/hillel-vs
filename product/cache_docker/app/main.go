package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {

	r := gin.Default()

	r.Static("/static", "./front/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"version": fmt.Sprintf("%v", time.Now().Unix()),
			"title":   "Title of site",
		})
	})

	v1 := r.Group("/api/v1")
	{
		//v1.POST("/register", registerHandler)
		//
		//v1.POST("/authorization", authorizationHandler)
		//
		//v1.POST("/delete", deleteHandler)

		v1.GET("/alert", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello, everyone!")
		})
		v1.POST("/post", postHandler)

		v1.GET("/getSessionID", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%v", time.Now().Unix()))
		})
	}

	//
	//POST /api/v1/register - записує користувача і в базу, і в кеш
	//POST /api/v1/authorization - для перевірки пароля дістає користувача з кеша, якщо його нема  в кеші - дістає з бази
	//POST /api/v1/delete
	//

	// Запуск сервера
	r.Run(":8045")
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
