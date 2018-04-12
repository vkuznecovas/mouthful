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
	AwsAccessKeyID            *string `json:"awsAccessKeyID,omitempty"`
	AwsSecretAccessKey        *string `json:"awsSecretAccessKey,omitempty"`
	AwsRegion                 *string `json:"awsRegion,omitempty"`
	SSLEnabled                *bool   `json:"sslEnabled,omitempty"`
}

// Moderation - moderation section of our config
type Moderation struct {
	Enabled                bool   `json:"enabled"`
	AdminPassword          string `json:"adminPassword"`
	SessionDurationSeconds int    `json:"sessionDurationSeconds"`
	MaxCommentLength       *int   `json:"maxCommentLength,omitempty"`
}

// Config - root of our config
type Config struct {
	Database   Database   `json:"database"`
	Honeypot   bool       `json:"honeypot"`
	Moderation Moderation `json:"moderation"`
	Client     Client     `json:"client"`
	API        API        `json:"api"`
}

// API - api configuration part
type API struct {
	Port         *int         `json:"port,omitempty"`
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
