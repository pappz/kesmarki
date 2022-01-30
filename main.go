package main

import (
	"github.com/webkeydev/logger"

	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	log = logger.NewLogger("KESMARKI")
	wg  sync.WaitGroup
)

func init() {
	logger.SetTxtLogger()
	osSigs := make(chan os.Signal, 1)
	signal.Notify(osSigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-osSigs
		log.Println("interrupt...")
		wg.Done()
	}()
}

func main() {
	wg.Add(1)
	err := createMQTTServer()
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}

	wg.Wait()
}
