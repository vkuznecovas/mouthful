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
}

// Config - root of our config
type Config struct {
	Database   Database
	Honeypot   bool
	Moderation Moderation
}
