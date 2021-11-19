package peer

import (
	"fmt"

	"github.com/hienduyph/genesis/database"
)

type PeerNode struct {
	IP          string           `json:"ip"`
	Port        uint64           `json:"port"`
	IsBootstrap bool             `json:"is_bootstrap"`
	IsActive    bool             `json:"is_active"`
	Connected   bool             `json:"connected"`
	Account     database.Account `json:"account"`
}

func (p PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}
