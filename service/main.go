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

func handleMessage(msg string) {
	switch msg {
	case "up":
		log.Printf("shutter up")
		shutterControl.Up()
	case "stop":
		log.Printf("shutter stop")
		shutterControl.Stop()
	case "down":
		log.Printf("shutter down")
		shutterControl.Down()
	}
}

func tearDown() {
	log.Printf("close broker")
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

	auth, err := service.NewFileAuth("users", "/etc/kesmarki/users")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	brokerService, err = service.NewBrokerService(auth)
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}

	brokerService.AddMsgHandler("kesmarki/shutter", handleMessage)

	log.Printf("MQTT broker listening on: %s", service.TcpAddress)
	log.Printf("Webscoket listener on: %s", service.WsAddress)

	wg.Wait()
}
