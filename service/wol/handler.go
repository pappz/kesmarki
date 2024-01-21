package wol

import (
	"github.com/pappz/kesmarki/mqtt"
	log "github.com/sirupsen/logrus"
)

type wolHandler struct {
	wol *BudafokiWol
}

func RegisterWolHandler(broker mqtt.BrokerService, wol *BudafokiWol) {
	s := wolHandler{wol: wol}

	broker.AddMsgHandler("kesmarki/wol/budafoki", s.handleClientMessage)
}

func (s wolHandler) handleClientMessage(msg string) {
	switch msg {
	case "up":
		log.Infof("budafoki wake up")
		s.wol.Wake()
	case "down":
		log.Errorf("invalid shutter switch action: %s", msg)
	}
}
