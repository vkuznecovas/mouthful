package global_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/global"
)

func TestCleanupTypeString(t *testing.T) {
	somethingElseEntirely := global.CleanupType(-1)

	assert.Equal(t, "CleanupType(-1)", somethingElseEntirely.String())
	assert.Equal(t, "Unconfirmed", global.Unconfirmed.String())
	assert.Equal(t, "Deleted", global.Deleted.String())
}
