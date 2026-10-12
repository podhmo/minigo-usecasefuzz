package source

// Recursive generics cannot be converted element-wise: the
// instantiated spec of Node[T] contains the instantiation itself
// ([]Node[T] -> []Node[int]), so inline expansion never terminates.
// The walk bails on the repeated pair, keeps the honest raw
// assignment, and the build fails loudly.
type Node[T any] []Node[T]

type A struct {
	V Node[int]
}
