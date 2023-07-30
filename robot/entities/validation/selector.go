package validation

import "strings"

func IsXpath(s string) bool {
	return strings.Contains(s, "//") || strings.Contains(s, "@")
}
