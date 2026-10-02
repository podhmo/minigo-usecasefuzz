package main

import (
	"fmt"
	"net/http"
)

// Build a request object — no network needed.
func main() {
	req, err := http.NewRequest("GET", "https://example.com/api?q=1", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("X-Test", "yes")
	fmt.Println(req.Method, req.URL.Host, req.URL.Path, req.Header.Get("X-Test"))
}
