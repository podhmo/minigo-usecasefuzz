package main

import (
	"encoding/xml"
	"fmt"
)

// Stdlib but reflect-driven: init reaches unsafe.Pointer via
// internal/abi — same boundary as lim-yaml.
type Server struct {
	Host string `xml:"host"`
	Port int    `xml:"port"`
}

func main() {
	var s Server
	err := xml.Unmarshal([]byte(`<s><host>x</host><port>8080</port></s>`), &s)
	fmt.Println(s, err)
}
