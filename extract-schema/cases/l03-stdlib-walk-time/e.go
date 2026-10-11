package app

import "time"

type Event struct {
	At time.Time `json:"at"`
}
