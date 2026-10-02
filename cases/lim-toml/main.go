package main

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

// External module + reflect-driven decode.
func main() {
	var v map[string]any
	err := toml.Unmarshal([]byte("a = 1\nb = \"two\"\n"), &v)
	fmt.Println(v, err)
}
