package source

// Ghost is undefined — the field's type cannot be resolved.
type SrcItem struct{ V Ghost }

type SrcUser struct {
	ID    int64
	Name  string
	Items []SrcItem
}
