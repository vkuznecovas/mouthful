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
