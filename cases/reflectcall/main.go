package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// dispatch invokes a registered handler by name with reflect-built
// arguments — the RPC/router/plugin-registry job: functions live in a
// map, callers carry them by name, and reflect.Value.Call bridges them.

var registry = map[string]any{}

func register(name string, fn any) {
	registry[name] = fn
}

func dispatch(name string, args ...any) ([]any, error) {
	fn, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown command %q", name)
	}
	in := make([]reflect.Value, len(args))
	for i, a := range args {
		in[i] = reflect.ValueOf(a)
	}
	out := reflect.ValueOf(fn).Call(in)
	res := make([]any, len(out))
	for i, o := range out {
		res[i] = o.Interface()
	}
	return res, nil
}

func greet(name string) string {
	return "hello, " + name
}

func add(a, b int) int {
	return a + b
}

func parseQty(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("bad qty %q", s)
	}
	if n <= 0 {
		return 0, fmt.Errorf("qty must be positive: %d", n)
	}
	return n, nil
}

func main() {
	register("greet", greet)
	register("add", add)
	register("parse.qty", parseQty)

	for _, call := range []struct {
		name string
		args []any
	}{
		{"greet", []any{"world"}},
		{"add", []any{19, 23}},
		{"parse.qty", []any{"42"}},
		{"parse.qty", []any{"x"}},
		{"upper", []any{"hi"}},
	} {
		res, err := dispatch(call.name, call.args...)
		if err != nil {
			fmt.Printf("%s: error %v\n", call.name, err)
			continue
		}
		if e, ok := res[len(res)-1].(error); ok && e != nil {
			fmt.Printf("%s: error %v\n", call.name, e)
			continue
		}
		trimmed := res[:1]
		if len(res) > 1 {
			trimmed = res[:len(res)-1]
		}
		fmt.Printf("%s: %s\n", call.name, strings.Trim(fmt.Sprint(trimmed), "[]"))
	}
}
