package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Homework() {
	r := gin.Default()

	// Реєстрація
	r.POST("/api/v1/register", func(c *gin.Context) {
		var json struct {
			Name     string `json:"name"`
			Password int    `json:"password"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "Received",
			"name":     json.Name,
			"password": json.Password,
		})
	})

	// Переклад
	r.POST("/api/v1/translate", func(c *gin.Context) {
		translation := c.DefaultQuery("text", "Nothing to translate")
		//sourceLang := c.DefaultQuery("sourceLang", "auto")
		//targetLang := c.DefaultQuery("lang", "en")

		var expression = "mock answer"
		//expression, err := googletranslatefree.Translate(translation, sourceLang, targetLang)
		//if err != nil {
		//	log.Println("Translation error:", err)
		//	c.String(http.StatusInternalServerError, "Translation failed: %v", err)
		//	return
		//}

		c.String(http.StatusOK, "Your text for translation, %s.\n Translated expression: %s ", translation, expression)
	})

	// Розрахунок
	r.POST("/api/v1/calculate", calculateHandler)

	// Запуск сервера
	r.Run(":9090")
}

func calculateHandler(c *gin.Context) {
	firstNumStr := c.PostForm("firstNum")
	secondNumStr := c.PostForm("secondNum")
	action := c.PostForm("action")

	firstNum, err1 := strconv.ParseFloat(firstNumStr, 64)
	secondNum, err2 := strconv.ParseFloat(secondNumStr, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid numbers provided",
		})
		return
	}

	var result float64
	var err error

	switch action {
	case "+":
		result = firstNum + secondNum
	case "-":
		result = firstNum - secondNum
	case "*":
		result = firstNum * secondNum
	case "/":
		if secondNum != 0 {
			result = firstNum / secondNum
		} else {
			err = fmt.Errorf("division by zero")
		}
	default:
		err = fmt.Errorf("invalid action provided")
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"firstNum":  firstNum,
			"secondNum": secondNum,
			"action":    action,
			"result":    result,
		})
	}
}
