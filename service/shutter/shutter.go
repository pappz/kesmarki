package shutter

import (
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Control struct {
	pinUp   rpio.Pin
	pinStop rpio.Pin
	pinDown rpio.Pin
	mutex   sync.Mutex
}

func NewControl() (*Control, error) {
	s := &Control{
		pinUp:   rpio.Pin(27),
		pinStop: rpio.Pin(22),
		pinDown: rpio.Pin(17),
	}

	if err := rpio.Open(); err != nil {
		return s, err
	}

	s.pinUp.Input()
	s.pinStop.Input()
	s.pinDown.Input()
	return s, nil
}

func (s *Control) Release() {
	s.mutex.Lock()
	_ = rpio.Close()
	s.mutex.Unlock()
}

func (s *Control) Up() {
	s.push(&s.pinUp)
}

func (s *Control) Stop() {
	s.push(&s.pinStop)
}

func (s *Control) Down() {
	s.push(&s.pinDown)
}

func (s *Control) push(p *rpio.Pin) {
	s.mutex.Lock()

	p.Output()
	p.PullDown()
	time.AfterFunc(500*time.Millisecond, func() {
		p.Input()
		s.mutex.Unlock()
	})
}
