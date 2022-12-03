package flower

import (
	"encoding/json"
	"github.com/pappz/kesmarki/tradfi"

	log "github.com/sirupsen/logrus"

	"github.com/pappz/kesmarki/mqtt"
)

type flowerHandler struct {
	broker  mqtt.BrokerService
	control Control
}

func RegisterFlowerHandler(broker mqtt.BrokerService, control Control) {
	s := flowerHandler{
		broker:  broker,
		control: control,
	}
	broker.AddMsgHandler("zigbee2mqtt/switch/flower", s.handleSwitch)
}

func (f flowerHandler) handleSwitch(msg string) {
	data := tradfi.SwitchPayload{}
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		log.Errorf("failed to parse switch payload: %s", err.Error())
		return
	}

	switch data.Action {
	case "on":
		log.Debugln("flower led demo")
		err := f.control.PlayDemo()
		if err != nil {
			log.Debugln("failed to send out 'demo' cmd to flower led")
		}
	case "off":
		log.Debugln("flower led off")
		err := f.control.Off()
		if err != nil {
			log.Debugln("failed to send 'off' cmd to flower led")
		}
	default:
		log.Errorf("invalid flower switch action: %s", data.Action)
	}

}
