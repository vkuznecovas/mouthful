package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/api"
)

func TestNormalizePath(t *testing.T) {
	res := api.NormalizePath("/test/1/2/3")
	assert.Equal(t, "/test/1/2/3/", res)
	res = api.NormalizePath("/test/1/2/3/")
	assert.Equal(t, "/test/1/2/3/", res)
	res = api.NormalizePath("test/1/2/3/")
	assert.Equal(t, "/test/1/2/3/", res)
	res = api.NormalizePath("test/1/2/3")
	assert.Equal(t, "/test/1/2/3/", res)

	res = api.NormalizePath("/test/1/2/index.html")
	assert.Equal(t, "/test/1/2/index.html", res)

	res = api.NormalizePath("test/1/2/index.html")
	assert.Equal(t, "/test/1/2/index.html", res)
}

func TestShortenAuthor(t *testing.T) {
	// 104 chars
	input1 := "author author author author author author author author author author author author author author author"
	res := api.ShortenAuthor(input1, 50)
	assert.Equal(t, 50, len(res))
	assert.Equal(t, "author author author author author author autho...", res)

	// 6 chars
	input2 := "author"
	res = api.ShortenAuthor(input2, 50)
	assert.Equal(t, len(input2), len(res))
	assert.Equal(t, input2, res)

	// 2 chars
	input3 := "aa"
	res = api.ShortenAuthor(input3, 3)
	assert.Equal(t, len(input3), len(res))
	assert.Equal(t, input3, res)

	// 4 chars
	input4 := "aaaa"
	res = api.ShortenAuthor(input4, 3)
	assert.Equal(t, len(input4), len(res))
	assert.Equal(t, input4, res)

	// 5 chars
	input5 := "aaaaa"
	res = api.ShortenAuthor(input5, 4)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, "a...", res)
}
