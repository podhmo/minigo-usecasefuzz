package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// dump renders any value as text — the "inspect an unknown value" job a
// debug printer, test helper, or tiny marshaler does: kind dispatch,
// struct field walks, slice indexing, sorted map traversal.

type LineItem struct {
	SKU   string
	Qty   int
	Price float64
}

type Order struct {
	ID     string
	Items  []LineItem
	BillTo *Address
	Labels map[string]string
}

type Address struct {
	City string
	Zip  string
}

func dump(v reflect.Value) string {
	if !v.IsValid() {
		return "<nil>"
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return "nil"
		}
		return "&" + dump(v.Elem())
	case reflect.Interface:
		return dump(v.Elem())
	case reflect.Struct:
		t := v.Type()
		var b strings.Builder
		b.WriteString(t.Name() + "{")
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%s: %s", t.Field(i).Name, dump(v.Field(i)))
		}
		b.WriteString("}")
		return b.String()
	case reflect.Slice:
		var b strings.Builder
		b.WriteString("[")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(dump(v.Index(i)))
		}
		b.WriteString("]")
		return b.String()
	case reflect.Map:
		byKey := map[string]reflect.Value{}
		var keys []string
		for _, k := range v.MapKeys() {
			s := fmt.Sprint(k.Interface())
			keys = append(keys, s)
			byKey[s] = k
		}
		sort.Strings(keys)
		var b strings.Builder
		b.WriteString("{")
		for i, k := range keys {
			if i > 0 {
				b.WriteString(" ")
			}
			fmt.Fprintf(&b, "%s: %s", k, dump(v.MapIndex(byKey[k])))
		}
		b.WriteString("}")
		return b.String()
	case reflect.String:
		return fmt.Sprintf("%q", v.String())
	case reflect.Bool:
		return fmt.Sprint(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprint(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprint(v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Float())
	default:
		return fmt.Sprint(v.Interface())
	}
}

func main() {
	o := Order{
		ID: "A-1042",
		Items: []LineItem{
			{SKU: "book", Qty: 2, Price: 1200},
			{SKU: "pen", Qty: 5, Price: 150},
		},
		BillTo: &Address{City: "Tokyo", Zip: "100-0001"},
		Labels: map[string]string{"priority": "high", "gift": "yes"},
	}
	fmt.Println(dump(reflect.ValueOf(o)))
	fmt.Println(dump(reflect.ValueOf(&o).Elem().Field(2)))
	fmt.Println(dump(reflect.ValueOf([]int{3, 1, 4})))
	fmt.Println(dump(reflect.ValueOf(map[string]int{"b": 2, "a": 1})))
	var nothing *Address
	fmt.Println(dump(reflect.ValueOf(nothing)))
}
