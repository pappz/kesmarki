package logic

import (
	"encoding/json"
	"github.com/pappz/kesmarki/mqtt"
	"github.com/pappz/kesmarki/shutter"
	log "github.com/sirupsen/logrus"
)

type tradfiSwitchPayload struct {
	Action      string `json:"action"`
	Batter      uint16 `json:"batter"`
	LinkQuality uint16 `json:"linkquality"`
}

type shutterHandler struct {
	shutterControl *shutter.Control
}

func RegisterShutterHandler(broker mqtt.BrokerService, shutterControl *shutter.Control) {
	s := shutterHandler{shutterControl}

	broker.AddMsgHandler("kesmarki/shutter", s.handleAndroidMessage)
	broker.AddMsgHandler("zigbee2mqtt/switch/shutter", s.handleSwitch)
}

func (s shutterHandler) handleAndroidMessage(msg string) {
	switch msg {
	case "up":
		log.Printf("shutter up")
		s.shutterControl.Up()
	case "stop":
		log.Printf("shutter stop")
		s.shutterControl.Stop()
	case "down":
		log.Printf("shutter down")
		s.shutterControl.Down()
	}
}

func (s shutterHandler) handleSwitch(msg string) {
	data := tradfiSwitchPayload{}
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		log.Errorf("failed to parse switch payload: %s", err.Error())
		return
	}

	switch data.Action {
	case "on":
		log.Debugln("shutter switch up")
		s.shutterControl.Up()
	case "off":
		log.Debugln("shutter switch down")
		s.shutterControl.Down()
	default:
		log.Errorf("invalid shutter switch action: %s", data.Action)
	}

}
