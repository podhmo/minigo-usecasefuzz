package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// External-module dependency: YAML decode.
func main() {
	var v map[string]any
	err := yaml.Unmarshal([]byte("a: 1\nb: two\n"), &v)
	fmt.Println(v, err)
}
