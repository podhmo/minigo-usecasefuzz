package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Classic word-frequency count over embedded text — the script-style
// "analyze a document" job.
const text = `Go is an open source programming language that makes it easy to
build simple, reliable, and efficient software. Go is expressive, concise,
clean, and efficient. Its concurrency mechanisms make it easy to write
programs that get the most out of multicore and networked machines, while
its novel type system enables flexible and modular program construction.`

func main() {
	freq := map[string]int{}
	for _, w := range strings.Fields(text) {
		w = strings.ToLower(strings.TrimFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		}))
		if w != "" {
			freq[w]++
		}
	}
	words := make([]string, 0, len(freq))
	for w := range freq {
		words = append(words, w)
	}
	sort.Slice(words, func(i, j int) bool {
		if freq[words[i]] != freq[words[j]] {
			return freq[words[i]] > freq[words[j]]
		}
		return words[i] < words[j]
	})
	fmt.Printf("%d distinct words\n", len(words))
	for i, w := range words {
		if i >= 8 {
			break
		}
		fmt.Printf("%-12s %d\n", w, freq[w])
	}
}
