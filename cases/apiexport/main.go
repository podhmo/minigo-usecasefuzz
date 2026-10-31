// apiexport: export user records the way a small admin tool does —
// JSON lines on stdout via json.Encoder, structured JSON logs via
// slog, and a text/template summary over plain maps that the program
// keeps editing afterwards. The records carry unexported bookkeeping
// fields (password hash, cache) that must never reach the output.
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/template"
)

type Address struct {
	City string `json:"city"`
	Zip  string `json:",omitempty"`
}

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Tags     []string `json:",omitempty"`
	Age      int      `json:",string"`
	Address  *Address `json:"address,omitempty"`
	password string
	cache    map[string]string
	audit
}

type audit struct {
	CreatedBy string `json:"created_by"`
	touched   int
}

type Role int

func (r Role) MarshalText() ([]byte, error) {
	return []byte([]string{"guest", "member", "admin"}[r]), nil
}

func main() {
	users := []User{
		{ID: 1, Name: "Ann", Email: "ann@example.com", Age: 31, Tags: []string{"ops"},
			Address: &Address{City: "Tokyo", Zip: "100-0001"}, password: "x1",
			cache: map[string]string{"k": "v"}, audit: audit{CreatedBy: "seed", touched: 3}},
		{ID: 2, Name: "Bob", Age: 27, password: "x2", audit: audit{CreatedBy: "import"}},
	}

	// 1. JSON lines export
	enc := json.NewEncoder(os.Stdout)
	for _, u := range users {
		if err := enc.Encode(u); err != nil {
			fmt.Println("encode error:", err)
		}
	}

	// 2. roles keyed by a TextMarshaler, pretty-printed
	roles := map[Role][]int{2: {1}, 1: {2}}
	enc.SetIndent("", "  ")
	enc.Encode(roles)
	fmt.Printf("encoder type: %T\n", enc)

	// 3. structured logs (time dropped for a stable diff)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	logger.Info("exported", "count", len(users))
	logger.Warn("missing email", "id", users[1].ID)
	logger.Debug("detail", "user", users[0].Name)

	// 4. text summary over plain maps, then keep editing the rows
	rows := []map[string]any{}
	for _, u := range users {
		rows = append(rows, map[string]any{"name": u.Name, "tags": len(u.Tags)})
	}
	t := template.Must(template.New("s").Parse("{{range .}}{{.name}}:{{.tags}} {{end}}\n"))
	var b strings.Builder
	if err := t.Execute(&b, rows); err != nil {
		fmt.Println("template error:", err)
	}
	fmt.Print(b.String())
	rows[0]["tags"] = 9
	rows[1]["note"] = "edited after render"
	fmt.Println(rows[0]["tags"], rows[1]["note"], len(rows[1]))

	// 5. a cyclic graph is reported, not looped on
	type node struct {
		Name string
		Next *node
	}
	n := &node{Name: "a"}
	n.Next = n
	if err := enc.Encode(n); err != nil {
		fmt.Println("cycle:", err)
	}
}
