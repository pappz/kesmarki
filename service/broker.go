package main

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
)

const (
	tcpAddress = ":1883"
	wsAddress  = ":1882"
)

var (
	server *mqtt.Server
)

func createMQTTServer() error {
	mqttAuth, err := newMqttAuth()
	if err != nil {
		return err
	}

	server = mqtt.New()
	tcp := listeners.NewTCP("t1", tcpAddress)

	err = server.AddListener(tcp, &listeners.Config{
		Auth: mqttAuth,
	})
	if err != nil {
		return err
	}

	ws := listeners.NewWebsocket("ws1", wsAddress)
	err = server.AddListener(ws, &listeners.Config{
		Auth: mqttAuth,
	})
	if err != nil {
		return err
	}

	server.Events.OnMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		handleMessage(pk.TopicName, string(pk.Payload))
		return pk, nil
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Error(err)
		}
	}()
	log.Printf("MQTT broker listening on: %s", tcpAddress)
	log.Printf("Webscoket listener on: %s", wsAddress)

	return nil
}

func handleMessage(topic string, msg string) {
	if topic != "shutter" {
		return
	}

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
