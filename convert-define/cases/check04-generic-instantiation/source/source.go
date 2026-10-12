package source

// Element-wise conversion through generic instantiations: the spec
// behind List[int] is parametric ([]T), so reaching the element type
// requires substituting the type arguments into the spec ([]T -> []int).
// Both shapes are exercised here — a decl over an instantiation (V) and
// a written instantiation (W).
type List[T any] []T
type SrcList List[int]

type A struct {
	V SrcList
	W List[int]
}
