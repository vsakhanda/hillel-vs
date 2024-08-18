package main

import (
	"log"
	"os"

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

	// Підписка на повідомлення

	_, err = nc.Subscribe(UserSetEvent.ToString(), processSetUserStruct)
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe(UserUpdateEvent.ToString(), processUpdateUserStruct)
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe(UserGetEvent.ToString(), func(m *nats.Msg) {
		log.Printf("Отримано дані про сессію : %s", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe(UserDeleteEvent.ToString(), func(m *nats.Msg) {
		log.Printf("Користувача видалено: %s", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	// Очікування повідомлень
	select {}
}

func processSetUserStruct(m *nats.Msg) {
	log.Printf("Отримано дані користувача : %s", string(m.Data))

	m.Respond([]byte("Successfully proceed"))
}
func processUpdateUserStruct(m *nats.Msg) {
	log.Printf("Отримано оновлені дані : %s", string(m.Data))

	m.Respond([]byte("Successfully proceed"))
}
