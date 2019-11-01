package readingtag

import "strings"

// IncludesString checks in an array includes given string
func IncludesString(array []string, lookup string) bool {
	for _, item := range array {
		if strings.EqualFold(item, lookup) {
			return true
		}
	}

	return false
}
