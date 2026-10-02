package destination

type DstSub struct{ V int64 }

type DstWrap struct {
	A   [3]int64
	Sub [2]DstSub
}
