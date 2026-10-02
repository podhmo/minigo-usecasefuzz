package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Slugify titles for URLs — replace-space-with-dash text munging.
var (
	nonalnum = regexp.MustCompile(`[^a-z0-9]+`)
	accents  = strings.NewReplacer(
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"á", "a", "à", "a", "â", "a",
		"ö", "o", "ó", "o", "ô", "o",
		"ü", "u", "ú", "u", "ù", "u",
		"ç", "c", "ñ", "n",
	)
)

func slug(s string) string {
	s = accents.Replace(strings.ToLower(s))
	s = nonalnum.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func main() {
	titles := []string{
		"Hello, World!",
		"Café Münchën — a tour",
		"  Spaces  and   dashes ",
		"日本語タイトル",
		"100% coverage (finally)",
	}
	for _, t := range titles {
		fmt.Printf("%-30q -> %q\n", t, slug(t))
	}

	// bonus: rune-aware helpers used by normalizers
	for _, r := range "A1-日本" {
		fmt.Printf("%c letter=%v digit=%v\n", r, unicode.IsLetter(r), unicode.IsDigit(r))
	}
}
