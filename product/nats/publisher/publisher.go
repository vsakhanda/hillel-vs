package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	godotenv.Load()

	natsUrl := os.Getenv("NATS_URL")

	// Підключення до NATS сервера
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	session := fmt.Sprint(time.Now().Unix())
	data := []byte(session)

	// Запит без відповіді
	err = Publish(nc, UserSessionEvent.ToString(), data)
	if err != nil {
		log.Fatal("Publish error:", err)
	}

	dataStr := []byte(`{"name":"Alex"}`)

	// Відправка запиту та отримання відповіді
	msg, err := nc.Request(UserStructEvent.ToString(), dataStr, nats.DefaultTimeout)
	if err != nil {
		log.Fatal("Request error:", err)
	} else {
		log.Printf("Отримано відповідь: %s", string(msg.Data))
	}

	// jetStream
	js := jetStream(nc)
	jetStreamPublisher(js)
}

func Publish(nc *nats.Conn, subject string, data []byte) error {
	msg := nats.NewMsg(subject)
	msg.Data = data

	return nc.PublishMsg(msg)
}
