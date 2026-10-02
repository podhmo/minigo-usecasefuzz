package source

type SrcSub struct{ V int64 }

type SrcWrap struct {
	P    *SrcSub
	Name string
}
