package source

type SrcSub struct{ V int64 }

type SrcWrap struct {
	A   [3]int64
	Sub [2]SrcSub
}
