package service

import (
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
	"github.com/webkeydev/logger"
	"sync"
)

const (
	TcpAddress = ":1883"
	WsAddress  = ":1882"
)

var (
	log = logger.NewLogger("BROKER")
)

type BrokerService struct {
	server   *mqtt.Server
	handlers *sync.Map
}

func NewBrokerService(authController auth.Controller) (BrokerService, error) {
	bs := BrokerService{
		mqtt.New(),
		&sync.Map{},
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

func (bs *BrokerService) AddMsgHandler(topic string, handler Handler) {
	log.Printf("add topic listener:  %s", topic)
	bs.handlers.Store(topic, handler)
	log.Printf("map: %v", &bs.handlers)
}

func (bs *BrokerService) Close() {
	_ = bs.server.Close()
}

func (bs *BrokerService) onConnected(cl events.Client, pk events.Packet) {
	log.Printf("client connected: %s", cl.ID)
}

func (bs *BrokerService) onMessage(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
	pkx = pk
	v, ok := bs.handlers.Load(pk.TopicName)
	if !ok {
		return
	}

	h, ok := v.(Handler)
	if !ok {
		return
	}

	h(string(pk.Payload))
	return
}

func (bs *BrokerService) onDisconnected(cl events.Client, err error) {
	log.Printf("client disconnected: %s", cl.ID)
}
