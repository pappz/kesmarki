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
	server  *mqtt.Server
	handler func(string, string)
}

func NewBrokerService(authController auth.Controller, handler func(topic string, msg string)) (BrokerService, error) {
	bs := BrokerService{
		mqtt.New(),
		handler,
	}

	bs.server.Events.OnMessage = bs.onMessage
	bs.server.Events.OnConnect = bs.onConnected
	bs.server.Events.OnDisconnect = bs.onDisconnected

	tcp := listeners.NewTCP("t1", TcpAddress)
	err := bs.server.AddListener(tcp, &listeners.Config{
		Auth: authController,
	})
	if err != nil {
		return bs, err
	}

	ws := listeners.NewWebsocket("ws1", WsAddress)
	err = bs.server.AddListener(ws, &listeners.Config{
		Auth: authController,
	})
	if err != nil {
		return bs, err
	}

	go func() {
		err := bs.server.Serve()
		if err != nil {
			log.Error(err)
		}
	}()
	return bs, nil
}

func (bs *BrokerService) Close() {
	log.Printf("close MQTT broker")
	_ = bs.server.Close()
}

func (bs *BrokerService) onConnected(cl events.Client, pk events.Packet) {
	log.Printf("client connected: %s", cl.ID)
}

func (bs *BrokerService) onMessage(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
	bs.handler(pk.TopicName, string(pk.Payload))
	return pk, nil
}

func (bs *BrokerService) onDisconnected(cl events.Client, err error) {
	log.Printf("client disconnected: %s", cl.ID)
}
