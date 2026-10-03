package destination

type DstIDs []int64
type DstLeafList []DstLeaf
type DstScores map[string]int64
type DstPtr *DstLeaf

type DstLeaf struct{ V int64 }

type DstWrap struct {
	IDs    DstIDs
	Items  DstLeafList
	Scores DstScores
	Ptr    DstPtr
	RawIDs DstIDs
}
