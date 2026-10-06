// Command micro is a synthetic probe for interpreter-level performance
// regressions. It stresses the paths the realworld profile keeps flagging
// (frame setup, operand stack, method-set walks, member selection, boxing)
// without needing any target checkout. Output is a checksum line so an
// output diff between two builds is also a correctness signal.
package main

import (
	"fmt"
	"sort"
	"strings"
)

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

type Stringer interface{ String() string }

type point struct{ x, y int }

func (p point) String() string { return fmt.Sprintf("%d,%d", p.x, p.y) }

func main() {
	sum := 0

	// call/frame churn: recursive descent
	sum += fib(28)

	// interface dispatch: repeated dynamic method calls
	var s Stringer = point{x: 3, y: 4}
	for i := 0; i < 150000; i++ {
		sum += len(s.String())
	}

	// member selection + value boxing: plain field reads
	p := point{x: 7, y: 9}
	for i := 0; i < 200000; i++ {
		sum += p.x * p.y
	}

	// append + map traffic
	var xs []int
	for i := 0; i < 100000; i++ {
		xs = append(xs, i)
	}
	m := map[string]int{}
	for i, x := range xs {
		m[fmt.Sprintf("k%d", i)] = x
	}
	for _, x := range xs {
		sum += m[fmt.Sprintf("k%d", x)]
	}

	// a little sorting to cover the reflect-lite builtin paths
	strs := []string{"delta", "alpha", "charlie", "bravo"}
	sort.Strings(strs)
	sum += len(strings.Join(strs, ":"))

	fmt.Println("checksum:", sum)
}
