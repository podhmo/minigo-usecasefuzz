package main

import (
	"fmt"
	"regexp"
	"strings"
)

// regexp method-set coverage — find/replace/split on a compiled
// pattern, the standard "rewrite text by pattern" job.
func main() {
	re := regexp.MustCompile(`(\w+)=(\d+)`)
	s := "a=1 b=22 c=333"

	fmt.Println("find all:", re.FindAllString(s, -1))
	fmt.Println("submatch:", re.FindStringSubmatch(s))
	fmt.Println("replace:", re.ReplaceAllString(s, "[$2]"))
	fmt.Println("replace literal:", re.ReplaceAllLiteralString(s, "<X>"))

	upper := re.ReplaceAllStringFunc(s, strings.ToUpper)
	fmt.Println("func:", upper)

	fmt.Println("split:", re.Split(s, -1))
	fmt.Println("match:", re.MatchString("z=9"), re.MatchString("nope"))

	// indices into the source
	idx := re.FindStringIndex(s)
	fmt.Println("first at:", s[idx[0]:idx[1]])

	// Compile error path
	_, err := regexp.Compile("(")
	fmt.Println("compile err:", err != nil)
}
