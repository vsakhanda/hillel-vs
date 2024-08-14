package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

// User структура для зберігання інформації про користувача
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Підключення до NATS
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Створення нового роутера GIN
	router := gin.Default()

	// Маршрути
	router.POST("/api/v1/user/reg", userRegister)
	router.POST("/api/v1/user/auth", userAuth)
	router.GET("/api/v1/user/get", userGet)

	// Запуск сервера
	router.Run(":8080")
}

func userRegister(c *gin.Context) {
	// Отримання даних з запиту
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Відправка NATS-запиту до сервісу Auth
	response, err := nc.Request("auth.register", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Відправка відповіді клієнту
	c.JSON(http.StatusOK, decode(response.Data))
}

func userAuth(c *gin.Context) {
	// Отримання даних з запиту
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Відправка NATS-запиту до сервісу Auth
	response, err := nc.Request("auth.authenticate", encode(user), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Відправка відповіді клієнту
	c.JSON(http.StatusOK, decode(response.Data))
}

func userGet(c *gin.Context) {
	// Отримання даних з запиту
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Відправка NATS-запиту до сервісу Users
	response, err := nc.Request("users.get", []byte(username), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Відправка відповіді клієнту
	c.JSON(http.StatusOK, decode(response.Data))
}

// Функції для кодування та декодування даних
func encode(data interface{}) []byte {
	encoded, _ := json.Marshal(data)
	return encoded
}

func decode(data []byte) map[string]interface{} {
	var decoded map[string]interface{}
	_ = json.Unmarshal(data, &decoded)
	return decoded
}
