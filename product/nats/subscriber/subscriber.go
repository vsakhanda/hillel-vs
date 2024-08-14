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
	_, err = nc.Subscribe(UserSessionEvent.ToString(), func(m *nats.Msg) {
		log.Printf("Отримано сесію: %s", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe(UserStructEvent.ToString(), processUserStruct)
	if err != nil {
		log.Fatal(err)
	}

	// // jetStream
	// js := jetStream(nc)
	// go jetStreamSubscriber(js)

	// Очікування повідомлень
	select {}
}

func processUserStruct(m *nats.Msg) {
	log.Printf("Отримано структуру: %s", string(m.Data))

	m.Respond([]byte("Successfully proceed"))
}
