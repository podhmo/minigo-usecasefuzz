package app

import "net"

type Stamp struct {
	IP net.IP `json:"ip"`
}
