package model

// ClientConfig - config for client
type ClientConfig struct {
	MaxCommentLength *int
	Honeypot         bool
	UseDefaultStyle  bool
	CustomCSSPath    *string
	Moderation       bool
	APIPort          *int
	APIHost          string
}
