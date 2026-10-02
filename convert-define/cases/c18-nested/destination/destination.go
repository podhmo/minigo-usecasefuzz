package destination

type DstInner struct{ ID int64 }

type DstUser struct {
	Inner *DstInner
	Flat  int64
}
