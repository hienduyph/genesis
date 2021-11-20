package peer

import (
	"fmt"

	"github.com/hienduyph/genesis/database"
)

type PeerNode struct {
	// core info, must be unique
	IP   string `json:"ip"`
	Port uint64 `json:"port"`

	// other meta could change
	Account database.Account `json:"account"`

	// state
	IsBootstrap bool `json:"is_bootstrap"`
	IsActive    bool `json:"is_active"`
	Connected   bool `json:"connected"`
}

func (p PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}
