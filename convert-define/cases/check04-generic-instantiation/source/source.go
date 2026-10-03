package source

// Generic instantiations cannot be converted element-wise: the spec
// behind List[int] is parametric ([]T), so the element types are
// invisible without type checking. Both failure shapes are pinned
// here — a decl over an instantiation (V) warns; a written
// instantiation (W) falls into the leaf mismatch warning. Both keep
// the honest raw assignment and fail to build.
type List[T any] []T
type SrcList List[int]

type A struct {
	V SrcList
	W List[int]
}
