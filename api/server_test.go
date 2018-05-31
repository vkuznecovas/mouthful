package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/api"
	configModel "github.com/vkuznecovas/mouthful/config/model"
)

var fakeKey = "key"
var fakeSecret = "secret"
var fakeOrigin = "http://fake.origin"
var fakeAdminUsers = []string{"id"}
var someFakeOauthProviders = []configModel.OauthProvider{
	configModel.OauthProvider{
		Enabled:      true,
		Name:         "github",
		Key:          &fakeKey,
		Secret:       &fakeSecret,
		AdminUserIds: &fakeAdminUsers,
	},
}

var serverTestConfig = configModel.Config{
	Honeypot: false,
	Moderation: configModel.Moderation{
		Enabled:              true,
		AdminPassword:        "sssssssssssss",
		MaxCommentLength:     &maxCommentLength,
		DisablePasswordLogin: false,
		OAauthProviders:      &someFakeOauthProviders,
		OAuthCallbackOrigin:  &fakeOrigin,
	},
	API: configModel.API{
		Debug:   false,
		Logging: false,
		Cache: configModel.Cache{
			Enabled:           false,
			IntervalInSeconds: 1,
			ExpiryInSeconds:   2,
		},
		RateLimiting: configModel.RateLimiting{
			Enabled:   false,
			PostsHour: 2,
		},
	},
	Client: configModel.Client{
		UseDefaultStyle: true,
		PageSize:        10,
	},
}

func TestCheckModerationVariablesOK(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.OAauthProviders = nil
	err := api.CheckModerationVariables(&configCopy)
	assert.Nil(t, err)
}

func TestCheckModerationVariablesNoPassword(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.AdminPassword = ""
	configCopy.Moderation.DisablePasswordLogin = true
	configCopy.Moderation.OAauthProviders = nil
	err := api.CheckModerationVariables(&configCopy)
	assert.NotNil(t, err)
	assert.Equal(t, "You have moderation enabled with no enabled OAUTH providers or admin password functionality. You will not be able to login on the admin panel. Please check your configuration", err.Error())
}

func TestCheckModerationVariablesDefaultPassword(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.AdminPassword = "somepassword"
	configCopy.Moderation.OAauthProviders = nil
	err := api.CheckModerationVariables(&configCopy)
	assert.NotNil(t, err)
	assert.Equal(t, "Please change the config.Moderation.AdminPassword value in config. Do not leave the default there", err.Error())
}

func TestCheckModerationVariablesDefaultPasswordDisabled(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.AdminPassword = "somepassword"
	configCopy.Moderation.DisablePasswordLogin = true
	configCopy.Moderation.OAauthProviders = nil
	err := api.CheckModerationVariables(&configCopy)
	assert.NotNil(t, err)
	assert.Equal(t, "You have moderation enabled with no enabled OAUTH providers or admin password functionality. You will not be able to login on the admin panel. Please check your configuration", err.Error())
}

func TestCheckModerationVariablesNoOrigin(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.OAuthCallbackOrigin = nil
	err := api.CheckModerationVariables(&configCopy)
	assert.NotNil(t, err)
	assert.Equal(t, "Please provide a OAuthCallbackOrigin in the config moderation section. For more info, refer to the documentation on github", err.Error())
}

func TestCheckModerationVariablesDefaultPasswordDisabledNoErrorWithOAuth(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.AdminPassword = "somepassword"
	configCopy.Moderation.DisablePasswordLogin = true
	err := api.CheckModerationVariables(&configCopy)
	assert.Nil(t, err)
}

func TestCheckModerationVariablesEmptyPassword(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.AdminPassword = ""
	configCopy.Moderation.DisablePasswordLogin = false
	err := api.CheckModerationVariables(&configCopy)
	assert.NotNil(t, err)
	assert.Equal(t, "config.Moderation.AdminPassword is not defined in config", err.Error())
}

func TestCheckModerationVariablesEmptyPasswordDisabledNoError(t *testing.T) {
	configCopy := serverTestConfig
	configCopy.Moderation.AdminPassword = ""
	configCopy.Moderation.DisablePasswordLogin = true
	err := api.CheckModerationVariables(&configCopy)
	assert.Nil(t, err)
}

func TestOriginGetsSuffixed(t *testing.T) {
	configCopy := serverTestConfig
	origin := "http://some.origin"
	configCopy.Moderation.OAuthCallbackOrigin = &origin
	err := api.CheckModerationVariables(&configCopy)
	assert.Nil(t, err)
	assert.Equal(t, "http://some.origin/", *configCopy.Moderation.OAuthCallbackOrigin)
}
