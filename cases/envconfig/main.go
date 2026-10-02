package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Environment-driven config: read env vars with typed fallbacks —
// the standard twelve-factor script pattern.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func main() {
	os.Setenv("APP_PORT", "9090")
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_TAGS", "a,b, c")

	fmt.Println("host:", getenv("APP_HOST", "127.0.0.1"))
	fmt.Println("port:", getenvInt("APP_PORT", 8080))
	fmt.Println("missing port:", getenvInt("APP_NOPE", 1234))
	fmt.Println("debug:", getenv("APP_DEBUG", "false") == "true")
	for _, t := range strings.Split(getenv("APP_TAGS", ""), ",") {
		fmt.Println("tag:", strings.TrimSpace(t))
	}
	fmt.Println("environ>0:", len(os.Environ()) > 0)
}
