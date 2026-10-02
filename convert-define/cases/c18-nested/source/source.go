package source

type SrcInner struct{ V int64 }

type SrcUser struct {
	ID int64
	In SrcInner
}
