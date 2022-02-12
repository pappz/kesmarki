package main

import (
	"github.com/webkeydev/logger"

	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	log     = logger.NewLogger("KESMARKI")
	wg      sync.WaitGroup
	shutter *Shutter
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

func tearDown() {
	mqttClose()
	shutter.Release()

	wg.Done()
	log.Println("bye")
}

func main() {
	var err error

	shutter, err = NewShutter()
	if err != nil {
		log.Fatalf("%s", err)
	}

	wg.Add(1)
	err = createMQTTServer()
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}

	wg.Wait()
}
