package source

// LegacyField was removed from SrcUser; the stale generated.go still
// references it, so package gen currently does not compile.
type SrcUser struct {
	ID   int64
	Name string
}
