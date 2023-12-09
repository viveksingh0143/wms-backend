package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	reg, err := regexp.Compile("[^a-z0-9-]+")
	if err != nil {
		panic(err)
	}
	text = reg.ReplaceAllString(text, "")

	reg, err = regexp.Compile("--+")
	if err != nil {
		panic(err)
	}
	text = reg.ReplaceAllString(text, "-")
	return text
}
