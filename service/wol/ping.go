package wol

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// minPingInterval suppresses bursts of probe requests: if a ping happened
// within this window, the cached result is returned instead of running a new
// probe.
const minPingInterval = 5 * time.Second

// pinger runs the system `ping` binary against a host and caches the last
// result so that events arriving in quick succession do not trigger repeated
// probes.
type pinger struct {
	host string

	mu     sync.Mutex
	last   time.Time
	online bool
	primed bool
}

func newPinger(host string) *pinger {
	return &pinger{host: host}
}

// Online reports whether the host is reachable. If a probe ran less than
// minPingInterval ago, the previous result is returned without pinging again.
func (p *pinger) Online() bool {
	p.mu.Lock()
	if p.primed && time.Since(p.last) < minPingInterval {
		online := p.online
		age := time.Since(p.last)
		p.mu.Unlock()
		log.Debugf("ping %s: using cached result online=%t (age %s < %s)", p.host, online, age.Round(time.Millisecond), minPingInterval)
		return online
	}
	p.mu.Unlock()

	online := p.probe()

	p.mu.Lock()
	p.online = online
	p.last = time.Now()
	p.primed = true
	p.mu.Unlock()

	return online
}

// probe runs a single ICMP echo via the system ping binary and returns true on
// a successful reply.
func (p *pinger) probe() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// -c 1: one echo request, -W 2: wait up to 2s for the reply.
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", "2", p.host)
	start := time.Now()
	out, err := cmd.CombinedOutput()
	took := time.Since(start).Round(time.Millisecond)

	online := err == nil
	if online {
		log.Debugf("ping %s: online (took %s)", p.host, took)
	} else {
		log.Debugf("ping %s: offline (took %s): %s; output: %s", p.host, took, err, strings.TrimSpace(string(out)))
	}
	return online
}
