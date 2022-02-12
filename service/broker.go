package main

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
)

var (
	server *mqtt.Server
)

func createMQTTServer() error {
	server = mqtt.New()
	address := ":1883"
	tcp := listeners.NewTCP("t1", address)

	err := server.AddListener(tcp, nil)
	if err != nil {
		return err
	}

	ws := listeners.NewWebsocket("t2", ":1882")
	err = server.AddListener(ws, nil)
	if err != nil {
		return err
	}

	server.Events.OnMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		handleMessage(string(pk.Payload))
		return pk, nil
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Error(err)
		}
	}()
	log.Printf("MQTT broker listening on: %s", address)
	log.Printf("Webscoket listener on: :1882")

	return nil
}

func handleMessage(msg string) {
	log.Printf("on message: %s", msg)
	switch msg {
	case "up":
		shutter.Up()
	case "stop":
		shutter.Stop()
	case "down":
		shutter.Down()
	}
}

func mqttClose() {
	log.Printf("close MQTT broker")
	_ = server.Close()
}
