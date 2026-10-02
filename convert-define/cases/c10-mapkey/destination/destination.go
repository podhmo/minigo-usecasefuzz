package destination

type DstKey struct{ K string }
type DstSub struct{ V int64 }

type DstWrap struct {
	M map[int64]DstSub
	N map[DstKey]DstSub
}
