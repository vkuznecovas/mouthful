package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/global"
)

func TestParseConfig_Returns_Error_On_Invalid_Input(t *testing.T) {
	res, err := config.ParseConfig([]byte("this is not a json"))
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestParseConfig_Returns_Config(t *testing.T) {
	res, err := config.ParseConfig([]byte(`{"moderation":{"enabled":true}}`))
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Moderation.Enabled)
}

func TestTransformConfigToClientConfig_Sets_Defaults(t *testing.T) {
	cfg := model.Config{}
	cfg.Honeypot = false
	cfg.Client.PageSize = 10
	cfg.Client.UseDefaultStyle = true
	cfg.Moderation = model.Moderation{
		Enabled: true,
	}
	clientConfig := config.TransformConfigToClientConfig(&cfg)
	res, err := json.Marshal(clientConfig)
	fmt.Println(res)
	assert.Nil(t, err)
	assert.False(t, clientConfig.Honeypot)
	assert.True(t, clientConfig.Moderation)
	assert.True(t, clientConfig.UseDefaultStyle)
	assert.Equal(t, cfg.Client.PageSize, clientConfig.PageSize)
	assert.Equal(t, global.DefaultCommentLengthLimit, *clientConfig.MaxCommentLength)
	assert.Equal(t, global.DefaultAuthorLengthLimit, *clientConfig.MaxAuthorLength)
}

func TestTransformConfigToClientConfig_Overrides_Defaults(t *testing.T) {
	cfg := model.Config{}
	cfg.Honeypot = false
	cfg.Client.PageSize = 10
	cfg.Client.UseDefaultStyle = true
	maxLen := 15
	cfg.Moderation = model.Moderation{
		Enabled:          true,
		MaxAuthorLength:  &maxLen,
		MaxCommentLength: &maxLen,
	}
	clientConfig := config.TransformConfigToClientConfig(&cfg)
	assert.False(t, clientConfig.Honeypot)
	assert.True(t, clientConfig.Moderation)
	assert.True(t, clientConfig.UseDefaultStyle)
	assert.Equal(t, cfg.Client.PageSize, clientConfig.PageSize)
	assert.Equal(t, maxLen, *clientConfig.MaxCommentLength)
	assert.Equal(t, maxLen, *clientConfig.MaxAuthorLength)
}

func TestTransformToAdminConfig_Sets_Defaults(t *testing.T) {
	cfg := model.Config{}
	res := config.TransformToAdminConfig(&cfg)
	assert.Equal(t, "/", res.Path)
	assert.False(t, res.DisablePasswordLogin)
	assert.Len(t, *res.OauthProviders, 0)
}

func TestTransformToAdminConfig_Overrides_Defaults(t *testing.T) {
	cfg := model.Config{}
	secret := "secret"
	key := "key"
	path := "/test"
	adminUserIds := &[]string{"test"}
	providers := &[]model.OauthProvider{model.OauthProvider{
		Name:         "test",
		Enabled:      true,
		Key:          &key,
		Secret:       &secret,
		AdminUserIds: adminUserIds,
	}}
	cfg.Moderation = model.Moderation{
		OAauthProviders:      providers,
		Path:                 &path,
		DisablePasswordLogin: true,
	}
	res := config.TransformToAdminConfig(&cfg)
	assert.True(t, res.DisablePasswordLogin)
	assert.Len(t, *res.OauthProviders, 1)
	cp := *res.OauthProviders
	assert.Equal(t, "test", cp[0])
	assert.Equal(t, path, res.Path)
}
