package source

type SrcKey struct{ K string }
type SrcSub struct{ V int64 }

type SrcWrap struct {
	M map[int]SrcSub
	N map[SrcKey]SrcSub
}
