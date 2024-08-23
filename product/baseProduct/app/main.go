package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

var (
	ctx             = context.Background()
	redisClient     *redis.Client
	mongoClient     *mongo.Client
	usersCollection *mongo.Collection
)

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func register(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword := hashPassword(body.Password)

	// Запис в MongoDB
	_, err := usersCollection.InsertOne(ctx, bson.M{"username": body.Username, "password": hashedPassword})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Запис в Redis
	err = redisClient.Set(ctx, body.Username, hashedPassword, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func authorization(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword := hashPassword(body.Password)

	// Перевірка в Redis
	cachedPassword, err := redisClient.Get(ctx, body.Username).Result()
	if err == redis.Nil {
		// Якщо користувача немає в Redis, шукаємо в MongoDB
		var user bson.M
		err := usersCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
			return
		}

		// Перевірка пароля
		if user["password"] != hashedPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
			return
		}

		// Кешування користувача в Redis
		redisClient.Set(ctx, body.Username, user["password"], 0)

		c.JSON(http.StatusOK, gin.H{"message": "Authorization successful"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check cache"})
		return
	}

	if cachedPassword != hashedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authorization successful"})
}

func delete(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Видалення з Redis
	err := redisClient.Del(ctx, body.Username).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from cache"})
		return
	}

	// Видалення з MongoDB
	_, err = usersCollection.DeleteOne(ctx, bson.M{"username": body.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
	// Підключення до Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Підключення до MongoDB
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		panic(err)
	}

	usersCollection = mongoClient.Database("user_db").Collection("users")

	// Налаштування роутера Gin
	router := gin.Default()

	router.POST("/api/v1/register", register)
	router.POST("/api/v1/authorization", authorization)
	router.POST("/api/v1/delete", delete)

	router.Run(":8045")
}
