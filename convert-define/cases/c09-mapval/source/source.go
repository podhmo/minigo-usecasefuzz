package source

type SrcSub struct{ V int64 }

type SrcWrap struct {
	M map[string]SrcSub
}
