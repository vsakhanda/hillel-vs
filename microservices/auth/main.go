package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

// User структура для зберігання інформації про користувача
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var sessions = map[string]string{} // Сесіі користувачів (ключ = userID, значення = sessionID)

func main() {
	// Підключення до NATS
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Підписка на теми
	nc.Subscribe("auth.register", handleRegister)
	nc.Subscribe("auth.authenticate", handleAuthenticate)

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

	// Відправка NATS-запиту до сервісу Users для реєстрації користувача
	response, err := nc.Request("users.register", m.Data, nats.DefaultTimeout)
	if err != nil {
		nc.Publish(m.Reply, []byte(fmt.Sprintf(`{"error": "Users service error: %v"}`, err)))
		return
	}

	// Відправка відповіді клієнту
	nc.Publish(m.Reply, response.Data)
}

func handleAuthenticate(m *nats.Msg) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(m.Data, &credentials)
	if err != nil {
		nc.Publish(m.Reply, []byte(fmt.Sprintf(`{"error": "Invalid data: %v"}`, err)))
		return
	}

	// Відправка NATS-запиту до сервісу Users для отримання інформації про користувача
	response, err := nc.Request("users.get", []byte(credentials.Username), nats.DefaultTimeout)
	if err != nil {
		nc.Publish(m.Reply, []byte(fmt.Sprintf(`{"error": "Users service error: %v"}`, err)))
		return
	}

	var user User
	err = json.Unmarshal(response.Data, &user)
	if err != nil || user.Password != credentials.Password {
		fmt.Printf("Unmarshal '%+v' '%x'\n", string(response.Data), err)
		fmt.Printf("handleAuthenticate '%s' '%s'\n", user.Password, credentials.Password)
		nc.Publish(m.Reply, []byte(`{"error": "Invalid username or password"}`))
		return
	}

	// Створення нової сесії
	sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
	sessions[user.ID] = sessionID

	// Відправка відповіді з сесією
	responseData, _ := json.Marshal(map[string]string{
		"status":    "Authenticated",
		"sessionID": sessionID,
	})
	nc.Publish(m.Reply, responseData)
}
