package main

import (
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Shutter struct {
	pinUp   rpio.Pin
	pinStop rpio.Pin
	pinDown rpio.Pin
	mutex   sync.Mutex
}

func NewShutter() (*Shutter, error) {
	s := &Shutter{
		pinUp:   rpio.Pin(23),
		pinStop: rpio.Pin(22),
		pinDown: rpio.Pin(24),
	}

	log.Printf("init gpio pins")
	if err := rpio.Open(); err != nil {
		return s, err
	}

	s.pinUp.Input()
	s.pinStop.Input()
	s.pinDown.Input()
	return s, nil
}

func (s *Shutter) Release() {
	log.Printf("release gpio resources")
	s.mutex.Lock()
	_ = rpio.Close()
	s.mutex.Unlock()
}

func (s *Shutter) Up() {
	s.push(&s.pinUp)
}

func (s *Shutter) Stop() {
	s.push(&s.pinStop)
}

func (s *Shutter) Down() {
	s.push(&s.pinDown)
}

func (s *Shutter) push(p *rpio.Pin) {
	s.mutex.Lock()

	p.Output()
	p.PullDown()
	time.AfterFunc(500*time.Millisecond, func() {
		p.Input()
		s.mutex.Unlock()
	})
}
