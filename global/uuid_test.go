package global_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/global"
)

func TestGeneratesValidUUID(t *testing.T) {
	uuid := global.GetUUID()
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	assert.True(t, r.MatchString(uuid.String()))
}

func TestParsesValidUUIDString(t *testing.T) {
	uuid := global.GetUUID()
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	assert.True(t, r.MatchString(uuid.String()))
	parsed, err := global.ParseUUIDFromString(uuid.String())
	assert.Nil(t, err)
	assert.True(t, r.MatchString(parsed.String()))
}

func TestParsesInvalidUUIDString(t *testing.T) {
	uuid := "535622c2-2da5-11e8-b4670ed5f89f718b"
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	assert.False(t, r.MatchString(uuid))
	parsed, err := global.ParseUUIDFromString(uuid)
	assert.NotNil(t, err)
	assert.Nil(t, parsed)
}
