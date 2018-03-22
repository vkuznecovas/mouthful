package model

// Database - a config object representing our database
type Database struct {
	Dialect  string
	Database string
	Username string
	Password string
	Host     string
	Port     string
}

// Moderation - moderation section of our config
type Moderation struct {
	Enabled                bool
	SessionSecret          string
	SessionName            *string
	AdminPassword          string
	SessionDurationSeconds int
	MaxCommentLength       *int
}

// Config - root of our config
type Config struct {
	Database   Database
	Honeypot   bool
	Moderation Moderation
	Client     Client
	API        API
}

// API - api configuration part
type API struct {
	StaticPath   *string
	Port         *int
	Host         string
	Debug        bool
	Cache        Cache
	RateLimiting RateLimiting
}

// Client - client configuration part
type Client struct {
	UseDefaultStyle bool
	PageSize        int
}

// RateLimiting - rate limiting configuration
type RateLimiting struct {
	Enabled   bool
	PostsHour int
}

// Cache - cache settings
type Cache struct {
	Enabled           bool
	ExpiryInSeconds   int
	IntervalInSeconds int
}
