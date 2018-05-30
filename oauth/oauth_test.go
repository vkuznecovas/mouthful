package oauth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/oauth"
)

var admins = []string{"test"}
var secret = "secret"
var key = "key"

var githubCfg = model.OauthProvider{
	Name:         "github",
	Secret:       &secret,
	Key:          &key,
	AdminUserIds: &admins,
	Enabled:      true,
}

var facebookCfg = model.OauthProvider{
	Name:         "facebook",
	Secret:       &secret,
	Key:          &key,
	AdminUserIds: &admins,
	Enabled:      true,
}

func TestGetProvidersReturnsEmptyArrayOfProvidersOnEmptyConfig(t *testing.T) {
	providers, err := oauth.GetProviders(nil, "t")
	assert.Nil(t, err)
	assert.Len(t, providers, 0)
}

func TestGetProvidersReturnsAnArrayOfProviders(t *testing.T) {
	inp := []model.OauthProvider{githubCfg, facebookCfg}
	providers, err := oauth.GetProviders(&inp, "/")
	assert.Nil(t, err)
	assert.Len(t, providers, 2)
}

func TestGetProvidersReturnsIgnoresDisabled(t *testing.T) {
	facebookCfg.Enabled = false
	defer func() { facebookCfg.Enabled = true }()
	inp := []model.OauthProvider{githubCfg, facebookCfg}
	providers, err := oauth.GetProviders(&inp, "/")
	assert.Nil(t, err)
	assert.Len(t, providers, 1)
}

func TestGetProvidersBadKey(t *testing.T) {
	facebookCfg.Enabled = false
	defer func() { facebookCfg.Enabled = true }()
	currentKey := githubCfg.Key
	defer func() { githubCfg.Key = currentKey }()
	githubCfg.Key = nil
	inp := []model.OauthProvider{githubCfg, facebookCfg}
	providers, err := oauth.GetProviders(&inp, "/")
	assert.NotNil(t, err)
	assert.Len(t, providers, 0)
	assert.Equal(t, "No key set for github OAUTH provider in config and environment variable GITHUB_KEY not set, cannot set up OAUTH for github", err.Error())
}
