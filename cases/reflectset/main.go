package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Config override merge: apply a flat string map onto a typed struct via
// reflection — the "env vars override the defaults" job.
type Server struct {
	Host  string `conf:"host"`
	Port  int    `conf:"port"`
	Debug bool   `conf:"debug"`
}

func applyOverrides(v any, kv map[string]string) error {
	rv := reflect.ValueOf(v).Elem()
	ty := rv.Type()
	for i := 0; i < ty.NumField(); i++ {
		name := ty.Field(i).Tag.Get("conf")
		raw, ok := kv[name]
		if !ok {
			continue
		}
		f := rv.Field(i)
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.String:
			f.SetString(raw)
		case reflect.Int:
			n, err := strconv.Atoi(raw)
			if err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			f.SetInt(int64(n))
		case reflect.Bool:
			f.SetBool(raw == "true" || raw == "1")
		}
	}
	return nil
}

func describe(v any) []string {
	rv := reflect.ValueOf(v)
	var out []string
	for i := 0; i < rv.NumField(); i++ {
		out = append(out, fmt.Sprintf("%s=%v", rv.Type().Field(i).Name, rv.Field(i).Interface()))
	}
	sort.Strings(out)
	return out
}

func main() {
	s := &Server{Host: "localhost", Port: 8080}
	err := applyOverrides(s, map[string]string{
		"port":  "9090",
		"debug": "true",
		"noop":  "ignored",
	})
	fmt.Println(describe(*s), err)

	bad := &Server{}
	err = applyOverrides(bad, map[string]string{"port": "not-a-number"})
	fmt.Println(err)
}
