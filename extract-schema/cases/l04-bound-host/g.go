package app

import "sync"

type Guard struct {
	Mu  sync.Mutex `json:"mu"`
	Once sync.Once `json:"once"`
}
