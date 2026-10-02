package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Worker-pool fan-out/fan-in — the canonical concurrent pipeline.
func main() {
	jobs := make(chan string)
	results := make(chan string)
	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		for j := range jobs {
			// worker id is nondeterministic across engines — report the
			// payload only so the oracle diff is stable.
			results <- strings.ToUpper(j)
		}
	}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i)
	}
	go func() {
		for _, j := range []string{"alpha", "beta", "gamma", "delta", "epsilon"} {
			jobs <- j
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	got := []string{}
	for r := range results {
		got = append(got, r)
	}
	sort.Strings(got)
	for _, g := range got {
		fmt.Println(g)
	}
}
