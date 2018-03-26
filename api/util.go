package api

import "strings"

// NormalizePath adds a missing slash at the front or the end of given input path
func NormalizePath(input string) string {
	if !strings.HasPrefix(input, "/") {
		input = "/" + input
	}
	if strings.LastIndex(input, ".") == -1 {
		if !strings.HasSuffix(input, "/") {
			input += "/"
		}
	}
	return input
}
