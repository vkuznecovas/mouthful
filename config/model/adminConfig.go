package model

// AdminConfig - config for admin
type AdminConfig struct {
	DisablePasswordLogin bool      `json:"disablePasswordLogin"`
	OauthProviders       *[]string `json:"oauthProviders,omitempty"`
	Path                 string    `json:"path"`
}
