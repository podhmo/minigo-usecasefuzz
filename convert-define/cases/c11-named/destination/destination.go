package destination

type DstStatus string

type DstWrap struct {
	S DstStatus
	N int64
}
