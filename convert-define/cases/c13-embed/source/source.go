package source

type SrcBase struct{ ID int64 }

type SrcUser struct {
	SrcBase
	Name string
}
