package main

import (
	"fmt"
	"sort"
	"strings"
)

// Parse INI text into sections and flatten to sorted key=value — the
// "read a .ini/.cfg" script job.
const ini = `# comment
[server]
host = example.com
port = 8080

[log]
level = debug
; another comment
file=/var/log/app.log

[toplevel]
bare=yes`

func main() {
	kv := map[string]string{}
	section := ""
	bad := 0
	for _, line := range strings.Split(ini, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			bad++
			continue
		}
		key := section + "." + strings.TrimSpace(k)
		kv[key] = strings.TrimSpace(v)
	}
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, "=", kv[k])
	}
	fmt.Println("bad:", bad)
}
