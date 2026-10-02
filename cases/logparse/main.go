package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Parse an embedded nginx-ish access log: counts per status code,
// requests per path, and a slow-request listing.
const logData = `192.168.0.1 - - [10/Oct/2025:13:55:36 +0000] "GET /api/users?page=2 HTTP/1.1" 200 1234 0.042
10.0.0.5 - - [10/Oct/2025:13:55:40 +0000] "POST /api/login HTTP/1.1" 401 89 0.011
192.168.0.1 - - [10/Oct/2025:13:56:01 +0000] "GET /static/logo.png HTTP/1.1" 200 5321 0.004
10.0.0.7 - - [10/Oct/2025:13:56:12 +0000] "GET /api/users HTTP/1.1" 200 2871 0.120
10.0.0.5 - - [10/Oct/2025:13:57:00 +0000] "GET /api/report HTTP/1.1" 500 0 1.845
172.16.0.3 - - [10/Oct/2025:13:57:33 +0000] "DELETE /api/users/42 HTTP/1.1" 204 0 0.067
10.0.0.7 - - [10/Oct/2025:13:58:21 +0000] "GET /missing HTTP/1.1" 404 153 0.002
malformed line that should not match
192.168.0.9 - - [10/Oct/2025:13:59:59 +0000] "GET /api/users?page=3 HTTP/1.1" 200 1190 0.038`

var re = regexp.MustCompile(`"([A-Z]+) ([^ ]+) HTTP/[0-9.]+" ([0-9]+) ([0-9]+) ([0-9.]+)`)

func main() {
	status := map[string]int{}
	paths := map[string]int{}
	slow := []string{}
	unmatched := 0
	for _, line := range strings.Split(logData, "\n") {
		m := re.FindStringSubmatch(line)
		if m == nil {
			unmatched++
			continue
		}
		method, path, code := m[1], m[2], m[3]
		status[code]++
		// strip query string for the path rollup
		if i := strings.Index(path, "?"); i >= 0 {
			path = path[:i]
		}
		paths[method+" "+path]++
		if secs, _ := strconv.ParseFloat(m[5], 64); secs >= 1.0 {
			slow = append(slow, fmt.Sprintf("%s %s -> %s in %ss", method, path, code, m[5]))
		}
	}
	codes := make([]string, 0, len(status))
	for c := range status {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	for _, c := range codes {
		fmt.Printf("status %s: %d\n", c, status[c])
	}
	ps := make([]string, 0, len(paths))
	for p := range paths {
		ps = append(ps, p)
	}
	sort.Slice(ps, func(i, j int) bool {
		if paths[ps[i]] != paths[ps[j]] {
			return paths[ps[i]] > paths[ps[j]]
		}
		return ps[i] < ps[j]
	})
	fmt.Println("-- top paths --")
	for i, p := range ps {
		if i >= 3 {
			break
		}
		fmt.Printf("%-24s %d\n", p, paths[p])
	}
	sort.Strings(slow)
	for _, s := range slow {
		fmt.Println("slow:", s)
	}
	fmt.Printf("unmatched lines: %d\n", unmatched)
}
