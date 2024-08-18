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

	//Запити з відповідями
	dataStr := []byte(`{"name":"Vlad", "age": "37", "user_id": "1" }`)
	msg, err := nc.Request(UserSetEvent.ToString(), dataStr, nats.DefaultTimeout)
	if err != nil {
		log.Fatal("Get user error:", err)
	} else {
		log.Printf("Отримано інформацію про користувача: %s", string(msg.Data))
	}

	dataStr1 := []byte(`{"name":"Vlad", "age": "36", "user_id": "1" }`)
	msg, err = nc.Request(UserUpdateEvent.ToString(), dataStr1, nats.DefaultTimeout)
	if err != nil {
		log.Fatal("Get user error:", err)
	} else {
		log.Printf("Оновлено інформацію про користувача: %s", string(msg.Data))
	}

	// Відправка запиту та отримання відповіді

	err = Publish(nc, UserGetEvent.ToString(), data)
	if err != nil {
		log.Fatal("Set user error:", err)
	}

	msg, err = nc.Request(UserDeleteEvent.ToString(), dataStr1, nats.DefaultTimeout)
	if err != nil {
		log.Fatal("Delete user error:", err)
	} else {
		log.Printf("Сесія користувача видалена %s", string(msg.Data))
	}

}

//
//func Publish(nc *nats.Conn, subject string, data []byte) error {
//	msg := nats.NewMsg(subject)
//	msg.Data = data
//
//	return nc.PublishMsg(msg)
//}

func Publish(nc *nats.Conn, subject string, data []byte) error {
	msg := nats.NewMsg(subject)
	msg.Data = data

	return nc.PublishMsg(msg)
}
