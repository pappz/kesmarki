package wol

import (
	"encoding/json"

	"github.com/pappz/kesmarki/mqtt"
	log "github.com/sirupsen/logrus"
)

const (
	cmdTopic       = "kesmarki/wol/budafoki"
	statusTopic    = "kesmarki/wol/budafoki/status"
	statusGetTopic = "kesmarki/wol/budafoki/status/get"
)

type statusMsg struct {
	Online bool `json:"online"`
}

type wolHandler struct {
	broker mqtt.BrokerService
	wol    *BudafokiWol
	pinger *pinger
}

func RegisterWolHandler(broker mqtt.BrokerService, wol *BudafokiWol) *wolHandler {
	s := &wolHandler{
		broker: broker,
		wol:    wol,
		pinger: newPinger(wol.IP()),
	}

	// Wake / control commands from the clients.
	broker.AddMsgHandler(cmdTopic, s.handleClientMessage)
	// Clients ask for a fresh status probe on connect.
	broker.AddMsgHandler(statusGetTopic, s.handleStatusRequest)

	return s
}

// PublishInitialStatus probes the machine once and publishes the retained
// status so it is populated before any client connects.
func (s *wolHandler) PublishInitialStatus() {
	s.publishStatus()
}

func (s *wolHandler) handleClientMessage(msg string) {
	switch msg {
	case "up":
		log.Infof("budafoki wake up")
		if err := s.wol.Wake(); err != nil {
			log.Errorf("failed to send WOL packet: %s", err)
		}
		s.publishStatus()
	case "down":
		log.Errorf("invalid shutter switch action: %s", msg)
	}
}

// handleStatusRequest is triggered when a client (e.g. the web app on connect)
// asks for the current machine state.
func (s *wolHandler) handleStatusRequest(_ string) {
	log.Debugln("budafoki status requested")
	s.publishStatus()
}

// publishStatus pings the machine (subject to the pinger's throttle) and
// publishes the result as a retained message so late-joining clients get the
// last known state immediately.
func (s *wolHandler) publishStatus() {
	online := s.pinger.Online()
	log.Infof("budafoki status: online=%t", online)
	payload, _ := json.Marshal(statusMsg{Online: online})
	if err := s.broker.Send(statusTopic, payload, true); err != nil {
		log.Errorf("failed to publish wol status: %s", err)
	}
}
