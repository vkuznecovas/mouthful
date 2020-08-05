package model

// Database - a config object representing our database
type Database struct {
	Dialect                   string  `json:"dialect"`
	Database                  *string `json:"database,omitempty"`
	Username                  *string `json:"username,omitempty"`
	Password                  *string `json:"password,omitempty"`
	Host                      *string `json:"host,omitempty"`
	Port                      *string `json:"port,omitempty"`
	TablePrefix               *string `json:"tablePrefix,omitempty"`
	DynamoDBThreadReadUnits   *int64  `json:"dynamoDBThreadReadUnits,omitempty"`
	DynamoDBCommentReadUnits  *int64  `json:"dynamoDBCommentReadUnits,omitempty"`
	DynamoDBThreadWriteUnits  *int64  `json:"dynamoDBThreadWriteUnits,omitempty"`
	DynamoDBCommentWriteUnits *int64  `json:"dynamoDBCommentWriteUnits,omitempty"`
	DynamoDBIndexWriteUnits   *int64  `json:"dynamoDBIndexWriteUnits,omitempty"`
	DynamoDBIndexReadUnits    *int64  `json:"dynamoDBIndexReadUnits,omitempty"`
	DynamoDBEndpoint          *string `json:"dynamoDBEndpoint,omitempty"`
	AwsAccessKeyID            *string `json:"awsAccessKeyID,omitempty"`
	AwsSecretAccessKey        *string `json:"awsSecretAccessKey,omitempty"`
	AwsRegion                 *string `json:"awsRegion,omitempty"`
	SSLEnabled                *bool   `json:"sslEnabled,omitempty"`
}

// Moderation - moderation section of our config
type Moderation struct {
	Enabled                bool             `json:"enabled"`
	AdminPassword          string           `json:"adminPassword"`
	DisablePasswordLogin   bool             `json:"disablePasswordLogin"`
	SessionDurationSeconds int              `json:"sessionDurationSeconds"`
	MaxCommentLength       *int             `json:"maxCommentLength,omitempty"`
	MaxAuthorLength        *int             `json:"maxAuthorLength,omitempty"`
	Path                   *string          `json:"path,omitempty"`
	OAauthProviders        *[]OauthProvider `json:"oauthProviders,omitempty"`
	OAuthCallbackOrigin    *string          `json:"oauthCallbackOrigin,omitempty"`
	PeriodicCleanUp        *PeriodicCleanUp `json:"periodicCleanup,omitempty"`
}

// Config - root of our config
type Config struct {
	Database     Database     `json:"database"`
	Honeypot     bool         `json:"honeypot"`
	Moderation   Moderation   `json:"moderation"`
	Client       Client       `json:"client"`
	API          API          `json:"api"`
	Notification Notification `json:"notification"`
}

// Notification - notification configuration part
type Notification struct {
	Webhook Webhook `json:"webhook"`
}

// Webhook represents the settings for notifications via webhook
type Webhook struct {
	Enabled bool    `json:"enabled"`
	URL     *string `json:"url"`
}

// API - api configuration part
type API struct {
	Port         *int         `json:"port,omitempty"`
	BindAddress  *string      `json:"bindAddress,omitempty"`
	Debug        bool         `json:"debug"`
	Cache        Cache        `json:"cache"`
	RateLimiting RateLimiting `json:"rateLimiting"`
	Cors         Cors         `json:"cors"`
	Logging      bool         `json:"logging"`
}

// Client - client configuration part
type Client struct {
	UseDefaultStyle bool `json:"useDefaultStyle"`
	PageSize        int  `json:"pageSize"`
}

// RateLimiting - rate limiting configuration
type RateLimiting struct {
	Enabled   bool `json:"enabled"`
	PostsHour int  `json:"postsHour"`
}

// Cache - cache settings
type Cache struct {
	Enabled           bool `json:"enabled"`
	ExpiryInSeconds   int  `json:"expiryInSeconds"`
	IntervalInSeconds int  `json:"entervalInSeconds"`
}

// Cors represents the cross origin resource sharing settings
type Cors struct {
	Enabled        bool      `json:"enabled"`
	AllowedOrigins *[]string `json:"allowedOrigins,omitempty"`
}

// OauthProvider represents the oauth provider configuration
type OauthProvider struct {
	Name         string    `json:"name"`
	Enabled      bool      `json:"enabled"`
	Secret       *string   `json:"secret,omitempty"`
	Key          *string   `json:"key,omitempty"`
	AdminUserIds *[]string `json:"adminUserIds,omitempty"`
}

// PeriodicCleanUp represents the settings for periodic cleanup of stale comments
type PeriodicCleanUp struct {
	Enabled                        bool  `json:"enabled"`
	RemoveDeleted                  bool  `json:"removeDeleted"`
	RemoveUnconfirmed              bool  `json:"removeUnconfirmed"`
	UnconfirmedTimeoutSeconds      int64 `json:"unconfirmedTimeoutSeconds"`
	DeletedTimeoutSeconds          int64 `json:"deletedTimeoutSeconds"`
	RemoveDeletedPeriodSeconds     int64 `josn:"removeDeletedPeriodSeconds"`
	RemoveUnconfirmedPeriodSeconds int64 `josn:"removeUnconfirmedPeriodSeconds"`
}
