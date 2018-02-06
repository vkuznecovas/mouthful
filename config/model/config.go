package model

type Database struct {
	Dialect  string
	Database string
	Username string
	Password string
	Host     string
	Port     string
}
type Moderation struct {
	Enabled                bool
	SessionSecret          string
	SessionName            *string
	AdminPassword          string
	SessionDurationSeconds int
}
type Config struct {
	Database   Database
	Honeypot   bool
	Moderation Moderation
}
