package destination

type DstSub struct{ V int64 }

type DstWrap struct {
	Items []DstSub
	Tags  []string
}
