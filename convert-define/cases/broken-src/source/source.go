package source

type SrcUser struct {
	ID   int64
	Name string
	Age  int
}

// CompileBroken does not compile (undefined identifier in the body) but
// parses fine — go/packages-based tools would refuse to run here.
func CompileBroken() int { return undefinedThing }
