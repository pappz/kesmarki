package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	formatter "github.com/webkeydev/logger"

	"github.com/pappz/kesmarki/ddns"
	"github.com/pappz/kesmarki/flower"
	"github.com/pappz/kesmarki/mqtt"
	"github.com/pappz/kesmarki/shutter"
	"github.com/pappz/kesmarki/wol"
)

var (
	wg             sync.WaitGroup
	shutterControl *shutter.Control
	brokerService  mqtt.BrokerService
	ddnsCancel     context.CancelFunc
)

type Config struct {
	WolBudafokiMac string

	DdnsHostedZoneID string
	DdnsRecordName   string
	DdnsInterval     time.Duration
}

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
	if ddnsCancel != nil {
		log.Printf("stop ddns updater")
		ddnsCancel()
	}
	log.Printf("close broker")
	brokerService.Close()
	log.Printf("release gpio resources")
	shutterControl.Release()

	wg.Done()
	log.Println("bye")
}

func main() {
	var err error

	cfg := readCfg()
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
	flower.RegisterFlowerHandler(brokerService, flowerControl)

	shutter.RegisterShutterHandler(brokerService, shutterControl)

	wolSrv := wol.NewBudafokiWol(cfg.WolBudafokiMac)
	wol.RegisterWolHandler(brokerService, wolSrv)

	startDdns(cfg)

	log.Printf("MQTT broker listening on: %s", mqtt.TcpAddress)
	log.Printf("Webscoket listener on: %s", mqtt.WsAddress)

	wg.Wait()
}

func startDdns(cfg Config) {
	if cfg.DdnsHostedZoneID == "" || cfg.DdnsRecordName == "" {
		log.Printf("ddns disabled (KM_DDNS_HOSTED_ZONE_ID / KM_DDNS_RECORD_NAME not set)")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	updater, err := ddns.New(ctx, ddns.Config{
		HostedZoneID: cfg.DdnsHostedZoneID,
		RecordName:   cfg.DdnsRecordName,
		Interval:     cfg.DdnsInterval,
	})
	if err != nil {
		cancel()
		log.Errorf("ddns init failed, continuing without it: %s", err)
		return
	}

	ddnsCancel = cancel
	go updater.Run(ctx)
}

func readCfg() Config {
	return Config{
		WolBudafokiMac: os.Getenv("KM_WOL_BUDAFOKI"),

		DdnsHostedZoneID: os.Getenv("KM_DDNS_HOSTED_ZONE_ID"),
		DdnsRecordName:   os.Getenv("KM_DDNS_RECORD_NAME"),
		DdnsInterval:     ddnsInterval(),
	}
}

func ddnsInterval() time.Duration {
	const fallback = time.Hour
	raw := os.Getenv("KM_DDNS_INTERVAL")
	if raw == "" {
		return fallback
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Warnf("invalid KM_DDNS_INTERVAL %q, using %s: %s", raw, fallback, err)
		return fallback
	}
	return d
}
