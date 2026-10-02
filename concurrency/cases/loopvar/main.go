package main

// Go 1.22+: each iteration has a fresh i, closures see 0,1,2.
func CapturePlain() int {
	fs := []func() int{}
	for i := 0; i < 3; i++ {
		fs = append(fs, func() int { return i })
	}
	return fs[0]() + fs[1]() + fs[2]() // want 3, shared-var => 9
}

func CaptureRangeInt() int {
	fs := []func() int{}
	for i := range 3 {
		fs = append(fs, func() int { return i })
	}
	return fs[0]() + fs[1]() + fs[2]() // want 3
}

func CaptureRangeSlice() int {
	fs := []func() int{}
	for _, v := range []int{10, 20, 30} {
		fs = append(fs, func() int { return v })
	}
	return fs[0]() + fs[1]() + fs[2]() // want 60, shared => 90
}

func CaptureIndexSlice() int {
	fs := []func() int{}
	for i, _ := range []int{10, 20, 30} {
		fs = append(fs, func() int { return i })
	}
	return fs[0]() + fs[1]() + fs[2]() // want 3, shared => 6
}

func main() {}
