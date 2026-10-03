package destination

type List[T any] []T
type DstList List[int64]

type B struct {
	V DstList
	W List[int64]
}
