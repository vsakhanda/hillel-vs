package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func jetStream(nc *nats.Conn) (js nats.JetStreamContext) {
	// Отримання контексту JetStream
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Створення потоку
	// _, err = js.AddStream(&nats.StreamConfig{
	// 	Name:     "ORDERS",
	// 	Subjects: []string{"orders.*"},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Потік створено")

	return js
}

func jetStreamSubscriber(js nats.JetStreamContext) {
	// Створення підписки на потік
	_, err := js.Subscribe("orders.*", func(m *nats.Msg) {
		log.Printf("Отримано повідомлення через jetStream, subject %s, data: %s", m.Subject, string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
}
