package main

import (
	"fmt"
	"reflect"
)

// cloneInto deep-copies src into dst via reflection — the "fork a config
// / take a snapshot" job: every level allocates through reflect.New and
// assigns through Set, so the copy shares nothing with the original.

type Feature struct {
	Name    string
	Enabled bool
}

type AppConfig struct {
	Name     string
	Replicas int
	Features []Feature
	Env      map[string]string
	Parent   *AppConfig
}

func cloneInto(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			cloneInto(dst.Field(i), src.Field(i))
		}
	case reflect.Slice:
		n := src.Len()
		cp := reflect.MakeSlice(src.Type(), n, n)
		for i := 0; i < n; i++ {
			cloneInto(cp.Index(i), src.Index(i))
		}
		dst.Set(cp)
	case reflect.Map:
		cp := reflect.MakeMap(src.Type())
		for _, k := range src.MapKeys() {
			e := reflect.New(src.Type().Elem()).Elem()
			cloneInto(e, src.MapIndex(k))
			cp.SetMapIndex(k, e)
		}
		dst.Set(cp)
	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		cp := reflect.New(src.Type().Elem())
		cloneInto(cp.Elem(), src.Elem())
		dst.Set(cp)
	default:
		dst.Set(src)
	}
}

func clone[T any](src T) T {
	dst := reflect.New(reflect.TypeOf(src)).Elem()
	cloneInto(dst, reflect.ValueOf(src))
	return dst.Interface().(T)
}

func main() {
	base := AppConfig{
		Name:     "api",
		Replicas: 3,
		Features: []Feature{{Name: "auth", Enabled: true}},
		Env:      map[string]string{"region": "tokyo"},
		Parent:   &AppConfig{Name: "shared", Replicas: 1},
	}

	fork := clone(base)
	fork.Replicas = 9
	fork.Features[0].Enabled = false
	fork.Env["region"] = "osaka"
	fork.Parent.Replicas = 7

	fmt.Println("base:", base.Replicas, base.Features[0].Enabled, base.Env["region"], base.Parent.Replicas)
	fmt.Println("fork:", fork.Replicas, fork.Features[0].Enabled, fork.Env["region"], fork.Parent.Replicas)

	empty := AppConfig{Parent: nil}
	got := clone(empty)
	fmt.Println("nil-parent:", got.Parent == nil)
}
