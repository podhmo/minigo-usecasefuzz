package destination

type DstSub struct{ V int64 }

type DstWrap struct {
	M map[string]DstSub
}
