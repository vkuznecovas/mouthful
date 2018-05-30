package provider_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/oauth/provider"
)

var secret = "secret"
var key = "key"

func TestProviderNoError(t *testing.T) {
	provider, err := provider.New("github", &secret, &key, []string{"asdasd"}, "/url")
	assert.Nil(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFallbackToEnvVarOnNoKey(t *testing.T) {
	envVal := os.Getenv("GITHUB_KEY")
	os.Setenv("GITHUB_KEY", "somekey")
	defer func() { os.Setenv("GITHUB_KEY", envVal) }()
	provider, err := provider.New("github", &secret, nil, []string{"asdasd"}, "/url")
	assert.Nil(t, err)
	assert.NotNil(t, provider)
}
func TestProviderFallbackToEnvVarOnNoSecret(t *testing.T) {
	envVal := os.Getenv("GITHUB_SECRET")
	os.Setenv("GITHUB_SECRET", "somekey")
	defer func() { os.Setenv("GITHUB_SECRET", envVal) }()
	provider, err := provider.New("github", nil, &key, []string{"asdasd"}, "/url")
	assert.Nil(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFallbackToEnvVarOnNone(t *testing.T) {
	envVal := os.Getenv("GITHUB_SECRET")
	os.Setenv("GITHUB_SECRET", "somekey")
	defer func() { os.Setenv("GITHUB_SECRET", envVal) }()
	envVal2 := os.Getenv("GITHUB_KEY")
	os.Setenv("GITHUB_KEY", "somekey")
	defer func() { os.Setenv("GITHUB_KEY", envVal2) }()
	provider, err := provider.New("github", nil, nil, []string{"asdasd"}, "/url")
	assert.Nil(t, err)
	assert.NotNil(t, provider)
}

func TestProviderNonExistant(t *testing.T) {
	_, err := provider.New("githubasdasd", &secret, &key, []string{"asdasd"}, "/url")
	assert.NotNil(t, err)
	assert.Equal(t, "No such OAUTH provider githubasdasd", err.Error())
}

func TestProviderNoSecret(t *testing.T) {
	_, err := provider.New("github", nil, &key, []string{"asdasd"}, "/url")
	assert.NotNil(t, err)
	assert.Equal(t, "No secret set for github OAUTH provider in config and environment variable GITHUB_SECRET not set, cannot set up OAUTH for github", err.Error())
}

func TestProviderNoKey(t *testing.T) {
	_, err := provider.New("github", &secret, nil, []string{"asdasd"}, "/url")
	assert.NotNil(t, err)
	assert.Equal(t, "No key set for github OAUTH provider in config and environment variable GITHUB_KEY not set, cannot set up OAUTH for github", err.Error())
}

func TestProviderNoAdmins(t *testing.T) {
	_, err := provider.New("github", &secret, &key, []string{}, "/url")
	assert.NotNil(t, err)
	assert.Equal(t, "No admin accounts provided for OAUTH provider github", err.Error())
}

func TestProviderInvalidUri(t *testing.T) {
	_, err := provider.New("github", &secret, &key, []string{"qq"}, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid callback uri provided for OAUTH provider github", err.Error())
}

func TestAllProvidersCreatedOk(t *testing.T) {
	for k := range provider.ProviderNameInitializationMap {
		p, err := provider.New(k, &secret, &key, []string{"qq"}, "/")
		assert.Nil(t, err)
		assert.NotNil(t, p)
	}
}
