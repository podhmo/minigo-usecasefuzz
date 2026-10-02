package main

import (
	"fmt"
	"net/url"
	"strings"
)

// Build and dissect a URL with escaping — link construction in a
// text pipeline.
func main() {
	base := "https://example.com"
	full, err := url.JoinPath(base, "docs", "v1", "spec.html")
	fmt.Println("join:", full, err)

	q := url.QueryEscape("hello world & more")
	p := url.PathEscape("a/b/c")
	fmt.Println("q:", q)
	fmt.Println("p:", p)

	raw := "https://example.com/search?q=" + q + "&tag=" + url.QueryEscape("go lang")
	fmt.Println("raw:", raw)

	// parse the query back out
	i := strings.Index(raw, "?")
	pairs := strings.Split(raw[i+1:], "&")
	for _, kv := range pairs {
		k, v, _ := strings.Cut(kv, "=")
		vv, err := url.QueryUnescape(v)
		fmt.Printf("%s=%v (err=%v)\n", k, vv, err)
	}
}
