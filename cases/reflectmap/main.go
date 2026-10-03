package main

import (
	"fmt"
	"reflect"
)

// Mini deserializer: fill a struct from map[string]any using field tags —
// the tiny "record → struct" mapper behind every codec.
type User struct {
	Name  string `db:"name"`
	Age   int    `db:"age"`
	Email string `db:"email,omitempty"`
	Admin bool   `db:"-"`
}

func fill(dst any, m map[string]any) {
	rv := reflect.ValueOf(dst).Elem()
	ty := rv.Type()
	for i := 0; i < ty.NumField(); i++ {
		sf := ty.Field(i)
		key := sf.Tag.Get("db")
		if key == "-" {
			continue
		}
		raw, ok := m[key]
		if !ok {
			continue
		}
		f := rv.Field(i)
		fv := reflect.ValueOf(raw)
		if fv.Type().AssignableTo(f.Type()) {
			f.Set(fv)
		}
	}
}

func main() {
	u := &User{}
	fill(u, map[string]any{
		"name":  "alice",
		"age":   30,
		"email": "a@example.com",
		"admin": true, // tag "-" — must be skipped
		"extra": "not a field",
	})
	fmt.Printf("%+v\n", *u)

	// type-identity check: a []User and []*User differ
	t1 := reflect.TypeOf([]User{})
	t2 := reflect.TypeOf([]*User{})
	fmt.Println(t1 == t2, t1.Kind(), t2.Elem().Kind())

	// value-of-value: inspecting a reflect.Value itself
	rv := reflect.ValueOf(reflect.ValueOf(42))
	fmt.Println(rv.Kind() == reflect.Struct, rv.Type().Name() != "")
}
