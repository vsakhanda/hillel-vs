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

	return js

}

func jetStreamPublisher(js nats.JetStreamContext) {
	// Публікація повідомлення у потік
	_, err := js.Publish("orders.created", []byte("Нове замовлення"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Повідомлення надіслане до потоку")
}
