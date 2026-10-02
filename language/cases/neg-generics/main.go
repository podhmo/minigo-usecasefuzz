package main

func Id[T any](x T) T { return x }
func Fst[T any](a, b T) T { return a }

func main() {
	_ = Id[int]("s")  // Go: cannot use "s" as int
	_ = Fst(1, "x")   // Go: cannot infer T
}
