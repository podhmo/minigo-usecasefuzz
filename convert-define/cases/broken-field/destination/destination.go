package destination

type DstItem struct{ V int64 }

type DstUser struct {
	ID    int64
	Name  string
	Items []DstItem
}
