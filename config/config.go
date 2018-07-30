// Package config deals with parsing and manipulation of the config files.
package config

import (
	"encoding/json"

	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/global"
)

// TODO: test

// ParseConfig takes in an absolute path for the config.json and uses it to create a config object.
func ParseConfig(contents []byte) (*model.Config, error) {
	var config model.Config
	err := json.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// TransformConfigToClientConfig transforms the config object to one that's safe to use for javascript by omitting all the sensitive fields.
func TransformConfigToClientConfig(input *model.Config) (conf *model.ClientConfig) {
	conf = &model.ClientConfig{}
	conf.Honeypot = input.Honeypot
	conf.Moderation = input.Moderation.Enabled
	if input.Moderation.MaxCommentLength != nil {
		conf.MaxCommentLength = input.Moderation.MaxCommentLength
	} else {
		length := global.DefaultCommentLengthLimit
		conf.MaxCommentLength = &length
	}
	conf.PageSize = input.Client.PageSize
	if input.Moderation.MaxAuthorLength != nil {
		conf.MaxAuthorLength = input.Moderation.MaxAuthorLength
	} else {
		length := global.DefaultAuthorLengthLimit
		conf.MaxAuthorLength = &length
	}
	conf.UseDefaultStyle = input.Client.UseDefaultStyle
	return conf
}

// TransformToAdminConfig takes in a config object and spits out an admin config
func TransformToAdminConfig(input *model.Config) (conf *model.AdminConfig) {
	providers := make([]string, 0)
	if input.Moderation.OAauthProviders != nil && len(*input.Moderation.OAauthProviders) > 0 {
		for _, v := range *input.Moderation.OAauthProviders {
			if v.Enabled {
				providers = append(providers, v.Name)
			}
		}
	}
	path := "/"
	if input.Moderation.Path != nil {
		path = *input.Moderation.Path
	}
	conf = &model.AdminConfig{
		DisablePasswordLogin: input.Moderation.DisablePasswordLogin,
		OauthProviders:       &providers,
		Path:                 path,
	}
	return conf
}
