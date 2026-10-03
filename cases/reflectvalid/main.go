package main

import (
	"fmt"
	"reflect"
	"sort"
)

// Struct-tag validation: `required:"true"` fields must be non-zero —
// the core of every request-validator package.
type Request struct {
	Path    string   `required:"true"`
	Method  string   `required:"true"`
	Retries int      `required:"false"`
	Headers []string `required:"true"`
}

func validate(v any) []string {
	rv := reflect.ValueOf(v)
	ty := rv.Type()
	var missing []string
	for i := 0; i < ty.NumField(); i++ {
		if ty.Field(i).Tag.Get("required") != "true" {
			continue
		}
		if rv.Field(i).IsZero() {
			missing = append(missing, ty.Field(i).Name)
		}
	}
	sort.Strings(missing)
	return missing
}

func main() {
	fmt.Println(validate(Request{Path: "/api", Method: "GET", Headers: []string{"x"}}))
	fmt.Println(validate(Request{Path: "/api"}))
	fmt.Println(validate(Request{}))

	// zero-vs-set distinction: an explicitly-zero int still fails required
	r := Request{Path: "/p", Method: "POST", Headers: []string{}}
	fmt.Println(validate(r)) // empty slice is non-nil but IsZero? no: allocated empty slice is not zero
}
