package wol

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

// BudafokiWol wakes up the budafoki PC over the network and can probe its
// current reachability by pinging its IP address.
type BudafokiWol struct {
	mac string
	ip  string
}

func NewBudafokiWol(mac, ip string) *BudafokiWol {
	return &BudafokiWol{
		mac: mac,
		ip:  ip,
	}
}

// IP returns the address that should be pinged to detect whether the machine
// is online.
func (w *BudafokiWol) IP() string {
	return w.ip
}

// Wake sends a Wake-on-LAN magic packet to the configured MAC address over a
// UDP broadcast on port 9.
func (w *BudafokiWol) Wake() error {
	packet, err := magicPacket(w.mac)
	if err != nil {
		return err
	}

	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return fmt.Errorf("dial broadcast: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("send magic packet: %w", err)
	}
	return nil
}

// magicPacket builds a Wake-on-LAN payload: 6 bytes of 0xFF followed by the
// target MAC repeated 16 times.
func magicPacket(mac string) ([]byte, error) {
	clean := strings.NewReplacer(":", "", "-", "").Replace(mac)
	hw, err := hex.DecodeString(clean)
	if err != nil || len(hw) != 6 {
		return nil, fmt.Errorf("invalid MAC address %q", mac)
	}

	packet := make([]byte, 0, 6+16*6)
	for i := 0; i < 6; i++ {
		packet = append(packet, 0xFF)
	}
	for i := 0; i < 16; i++ {
		packet = append(packet, hw...)
	}
	return packet, nil
}
