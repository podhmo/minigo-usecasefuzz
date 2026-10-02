package destination

import "example.com/m/box"

type DstWrap struct {
	Items []box.Box[int]
}
