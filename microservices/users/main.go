package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

// User структура для зберігання інформації про користувача
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	users = map[string]User{}
	mu    sync.Mutex
)

func main() {
	// Підключення до NATS
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Підписка на теми
	nc.Subscribe("users.register", handleRegister)
	nc.Subscribe("users.get", handleGetUser)

	// Тримання підключення активним
	select {}
}

func handleRegister(m *nats.Msg) {
	var user User
	err := json.Unmarshal(m.Data, &user)
	if err != nil {
		nc.Publish(m.Reply, []byte(fmt.Sprintf(`{"error": "Invalid data: %v"}`, err)))
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[user.Username]; exists {
		nc.Publish(m.Reply, []byte(`{"error": "User already exists"}`))
		return
	}

	user.ID = fmt.Sprintf("%d", len(users)+1)
	users[user.Username] = user

	responseData, _ := json.Marshal(map[string]string{
		"status": "User registered successfully",
		"userID": user.ID,
	})
	nc.Publish(m.Reply, responseData)
}

func handleGetUser(m *nats.Msg) {
	username := string(m.Data)

	mu.Lock()
	defer mu.Unlock()

	user, exists := users[username]
	if !exists {
		nc.Publish(m.Reply, []byte(`{"error": "User not found"}`))
		return
	}

	responseData, _ := json.Marshal(user)
	nc.Publish(m.Reply, responseData)
}
