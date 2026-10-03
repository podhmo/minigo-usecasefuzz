package main

import (
	"context"
	"fmt"
	"github.com/podhmo/minigo"
)

func main() {
	e := minigo.NewEngine("/private/tmp/minigo-reflect-mres/session-policy", minigo.WithPackageModes(map[string]minigo.PackageMode{"strings": minigo.ModeDeny}))
	_, err := e.Run(context.Background(), "./script", "Main")
	fmt.Println("parent denied:", err != nil)
	v, err := e.NewSession().Run(context.Background(), "./script", "Main")
	fmt.Println("session denied:", err != nil, "value:", v)
}
