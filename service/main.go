package main

import (
	"github.com/pappz/kesmarki/flower"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	formatter "github.com/webkeydev/logger"

	"github.com/pappz/kesmarki/mqtt"
	"github.com/pappz/kesmarki/shutter"
)

var (
	wg             sync.WaitGroup
	shutterControl *shutter.Control
	brokerService  mqtt.BrokerService
)

func init() {
	setLogFormatter()
	registerInterruptSignals()
}

func setLogFormatter() {
	formatter.SetTxtFormatterForLogger(log.StandardLogger())
	log.StandardLogger().SetLevel(log.DebugLevel)
}

func registerInterruptSignals() {
	osSigs := make(chan os.Signal, 1)
	signal.Notify(osSigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-osSigs
		log.Println("interrupt...")
		tearDown()
	}()
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
	log.Printf("start mqtt broker")

	auth, err := mqtt.NewFileAuth("users", "/etc/kesmarki/users")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	brokerService, err = mqtt.NewBrokerService(auth)
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}

	flowerControl := flower.NewControl(brokerService)

	shutter.RegisterShutterHandler(brokerService, shutterControl)
	flower.RegisterFlowerHandler(brokerService, flowerControl)

	log.Printf("MQTT broker listening on: %s", mqtt.TcpAddress)
	log.Printf("Webscoket listener on: %s", mqtt.WsAddress)

	wg.Wait()
}
