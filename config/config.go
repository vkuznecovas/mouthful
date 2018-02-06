package config

import (
	"encoding/json"

	"github.com/vkuznecovas/mouthful/config/model"
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
