package utils

import "regexp"

func HasFormatVerbs(format string) bool {
	re := regexp.MustCompile(`%[a-zA-Z]`)
	return re.MatchString(format)
}
