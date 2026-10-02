package source

type SrcSub struct{ V int64 }

type SrcWrap struct {
	Items []SrcSub
	Tags  []string
}
