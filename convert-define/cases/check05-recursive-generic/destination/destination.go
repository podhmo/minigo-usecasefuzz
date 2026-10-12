package destination

type Node[T any] []Node[T]

type B struct {
	V Node[int64]
}
