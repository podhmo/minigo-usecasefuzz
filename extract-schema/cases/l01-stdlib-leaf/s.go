package app

import (
	"encoding/json"
	"math/big"
	"net"
	"net/url"
	"time"
)

type Stamp struct {
	At   time.Time       `json:"at"`
	For  time.Duration   `json:"for"`
	IP   net.IP          `json:"ip"`
	URL  url.URL         `json:"url"`
	Big  big.Int         `json:"big"`
	Raw  json.RawMessage `json:"raw"`
}
