package source

// Named composites: the declared field types spell identifiers, so a
// conversion must unwrap the declaration (`type SrcIDs []int`) to find
// the element-wise shape — a plain DstT(src) cast cannot compile when
// the underlying types differ. See check03-named-slice for the
// still-failing sibling (incompatible element types).
type SrcIDs []int
type SrcLeafList []SrcLeaf
type SrcScores map[string]int
type SrcPtr *SrcLeaf

type SrcLeaf struct{ V int }

type SrcWrap struct {
	IDs    SrcIDs
	Items  SrcLeafList
	Scores SrcScores
	Ptr    SrcPtr
	RawIDs []int
}
