package service

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
	"github.com/webkeydev/logger"
)

const (
	TcpAddress = ":1883"
	WsAddress  = ":1882"
)

var (
	log = logger.NewLogger("BROKER")
)

type BrokerService struct {
	server *mqtt.Server
}

func NewBrokerService(authController auth.Controller, handler func(topic string, msg string)) (bs BrokerService, err error) {
	bs.server = mqtt.New()
	tcp := listeners.NewTCP("t1", TcpAddress)

	err = bs.server.AddListener(tcp, &listeners.Config{
		Auth: authController,
	})
	if err != nil {
		return
	}

	ws := listeners.NewWebsocket("ws1", WsAddress)
	err = bs.server.AddListener(ws, &listeners.Config{
		Auth: authController,
	})
	if err != nil {
		return
	}

	bs.server.Events.OnMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		handler(pk.TopicName, string(pk.Payload))
		return pk, nil
	}

	go func() {
		err := bs.server.Serve()
		if err != nil {
			log.Error(err)
		}
	}()
	return
}

func (bs *BrokerService) Close() {
	log.Printf("close MQTT broker")
	_ = bs.server.Close()
}
