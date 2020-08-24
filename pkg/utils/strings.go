package utils

import "strings"

// JoinNotEmptyStrings joins strings in `elem` using `sep` as a separator, ignoring the empty strings
func JoinNotEmptyStrings(sep string, elem ...string) string {
	tokens := make([]string, 0)
	for _, s := range elem {
		if s != "" {
			tokens = append(tokens, s)
		}
	}
	return strings.Join(tokens, sep)
}

// StringSliceIndexOf returns the index of an element in a string slice. If the slice doesn't contain the element, it will return -1
func StringSliceIndexOf(slice []string, elem string) int {
	for i, s := range slice {
		if s == elem {
			return i
		}
	}
	return -1
}
