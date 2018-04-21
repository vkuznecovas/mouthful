package global_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/global"
)

func TestSanitize(t *testing.T) {
	input := `hello <a name="n" href="javascript:alert('xss')">*you*</a>`
	res := global.ParseAndSaniziteMarkdown(input)
	assert.Equal(t, "<p>hello <em>you</em></p>\n", res)

	input = `[some text](javascript:alert('xss'))`
	res = global.ParseAndSaniziteMarkdown(input)
	assert.Equal(t, "<p><a title=\"xss\">some text</a>)</p>\n", res)

	input = `> hello <a name="n"
	> href="javascript:alert('xss')">*you*</a>`
	res = global.ParseAndSaniziteMarkdown(input)
	assert.Equal(t, "<blockquote>\n<p>hello  href=“javascript:alert(‘xss’)”&gt;<em>you</em></p>\n</blockquote>\n", res)

	input = `<script>alert("test")</script>`
	res = global.ParseAndSaniziteMarkdown(input)
	assert.Equal(t, "", res)
}
