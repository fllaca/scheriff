package utils

import "strings"

func JoinNotEmptyStrings(sep string, elem ...string) string {
	tokens := make([]string, 0)
	for _, s := range elem {
		if s != "" {
			tokens = append(tokens, s)
		}
	}
	return strings.Join(tokens, sep)
}
