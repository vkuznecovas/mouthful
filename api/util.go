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

// ShortenAuthor shortens the author name to an acceptable lenght, suffixing it with ...
func ShortenAuthor(input string, allowedLength int) string {
	// if the allowedLength is 3 or less, we can't really replace anything
	if allowedLength < 4 {
		return input
	}
	if len(input) > allowedLength {
		return input[0:allowedLength-3] + "..."
	}
	return input
}
