package global

import (
	bluemonday "github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// ParseAndSaniziteMarkdown takes in a markdown string, parses it and sanitizes it with blue monday.
// The returned html string should be safe for consumption.
func ParseAndSaniziteMarkdown(input string) string {
	unsafe := blackfriday.Run([]byte(input))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}
