package main

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
)

func createMQTTServer() error {
	server := mqtt.New()
	address := ":1883"
	tcp := listeners.NewTCP("t1", address)

	err := server.AddListener(tcp, nil)
	if err != nil {
		return err
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Error(err)
		}
	}()
	log.Printf("MQTT broker listening on: %s", address)
	return nil
}
