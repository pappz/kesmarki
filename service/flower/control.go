package flower

import (
	"encoding/json"
	"github.com/pappz/kesmarki/mqtt"
)

const (
	topic = "kesmarki/light/flower"
)

type msg struct {
	Action string `json:"action"`
}

type Control struct {
	broker mqtt.BrokerService
}

func NewControl(broker mqtt.BrokerService) Control {
	return Control{broker: broker}
}

func (c *Control) PlayDemo() error {
	m := msg{
		"demo",
	}
	payload, _ := json.Marshal(m)
	return c.broker.Send(topic, payload, true)
}

func (c *Control) Off() error {
	m := msg{
		"off",
	}
	payload, _ := json.Marshal(m)
	return c.broker.Send(topic, payload, true)
}
