// Package model deals with the model definitions of the configuration file.
package model

// ClientConfig - config for client
type ClientConfig struct {
	MaxCommentLength *int `json:"maxCommentLength,omitempty"`
	MaxAuthorLength  *int `json:"maxAuthorLength,omitempty"`
	Honeypot         bool `json:"honeypot"`
	UseDefaultStyle  bool `json:"useDefaultStyle"`
	Moderation       bool `json:"moderation"`
	PageSize         int  `json:"pageSize"`
}
