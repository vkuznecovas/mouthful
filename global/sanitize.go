package global

import (
	bluemonday "github.com/microcosm-cc/bluemonday"
	blackfriday "github.com/russross/blackfriday/v2"
)

// ParseAndSaniziteMarkdown takes in a markdown string, parses it and sanitizes it with blue monday.
// The returned html string should be safe for consumption.
func ParseAndSaniziteMarkdown(input string) string {
	unsafe := blackfriday.Run([]byte(input))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	htmlString := string(html)
	// in case we only got a <script> element or any other sort attempted injection, we'll get an empty paragraph
	// for future use, it's easier if we just return an empty string in that case, so we can validate against that
	if htmlString == "<p></p>\n" {
		return ""
	}
	return htmlString
}
