package main

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Load a JSON config into a typed struct plus a loose overlay map,
// merge, and re-emit as indented JSON — the "edit settings.json" job.
type Server struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	Tags []string `json:"tags,omitempty"`
}

type Config struct {
	Name    string            `json:"name"`
	Debug   bool              `json:"debug"`
	Server  Server            `json:"server"`
	Env     map[string]string `json:"env"`
	Weights map[string]int    `json:"weights"`
}

const raw = `{
	"name": "demo",
	"debug": false,
	"server": {"host": "localhost", "port": 8080},
	"env": {"HOME_DIR": "/srv"},
	"weights": {"a": 1}
}`

func main() {
	var cfg Config
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		fmt.Println("unmarshal:", err)
		return
	}
	cfg.Debug = true
	cfg.Server.Tags = []string{"api", "v2"}
	cfg.Env["LOG_LEVEL"] = "debug"
	cfg.Weights["b"] = 2

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("marshal:", err)
		return
	}
	fmt.Println(string(out))

	// loose access: pull a value back out generically
	var loose map[string]any
	if err := json.Unmarshal(out, &loose); err != nil {
		fmt.Println("loose:", err)
		return
	}
	srv := loose["server"].(map[string]any)
	fmt.Println("port as number:", srv["port"])
	keys := []string{}
	for k := range loose["env"].(map[string]any) {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("env keys:", keys)
}
