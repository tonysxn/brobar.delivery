package helpers

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var (
	specialCharsRegex = regexp.MustCompile(`[^\w\s-]`)
	multiDashRegex    = regexp.MustCompile(`-+`)
)

func GenerateSlug(input string) string {
	transliterated := unidecode.Unidecode(input)

	lower := strings.ToLower(transliterated)

	latinOnly := regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(lower, "-")

	withDashes := strings.ReplaceAll(latinOnly, " ", "-")

	singleDashes := regexp.MustCompile(`-+`).ReplaceAllString(withDashes, "-")

	trimmed := strings.Trim(singleDashes, "-")

	if trimmed == "" {
		return "untitled"
	}

	return trimmed
}
