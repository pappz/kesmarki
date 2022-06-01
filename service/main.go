package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pappz/kesmarki/service"
	"github.com/pappz/kesmarki/shutter"
	"github.com/webkeydev/logger"
)

var (
	log            = logger.NewLogger("KESMARKI")
	wg             sync.WaitGroup
	shutterControl *shutter.Control
	brokerService  service.BrokerService
)

func init() {
	logger.SetTxtLogger()
	osSigs := make(chan os.Signal, 1)
	signal.Notify(osSigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-osSigs
		log.Println("interrupt...")
		tearDown()
	}()
}

func handleMessage(topic string, msg string) {
	if topic != "shutter" {
		return
	}

	switch msg {
	case "up":
		shutterControl.Up()
	case "stop":
		shutterControl.Stop()
	case "down":
		shutterControl.Down()
	}
}

func tearDown() {
	brokerService.Close()
	log.Printf("release gpio resources")
	shutterControl.Release()

	wg.Done()
	log.Println("bye")
}

func main() {
	var err error

	log.Printf("init gpio pins")
	shutterControl, err = shutter.NewControl()
	if err != nil {
		log.Fatalf("%s", err)
	}

	wg.Add(1)
	log.Printf("start service broker")

	auth, err := service.NewFileAuth("/etc/kesmarki/users")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	brokerService, err = service.NewBrokerService(auth, handleMessage)
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}

	log.Printf("MQTT broker listening on: %s", service.TcpAddress)
	log.Printf("Webscoket listener on: %s", service.WsAddress)

	wg.Wait()
}
