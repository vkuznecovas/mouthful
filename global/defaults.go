// Package global provides constants, defaults and helper methods to mouthful.
package global

// StaticPath is the default location for our static files
const StaticPath = "./static"

// DefaultBindAddress is the default value for api service bind address
const DefaultBindAddress = "0.0.0.0"

// DefaultPort is the default value for api port
const DefaultPort = 8080

// DefaultDynamoDbThreadTableName default suffix for dynamodb thread
const DefaultDynamoDbThreadTableName = "mouthful_thread"

// DefaultDynamoDbCommentTableName default suffix for dynamodb comment
const DefaultDynamoDbCommentTableName = "mouthful_comment"

// DefaultCommentLengthLimit default comment length limit
const DefaultCommentLengthLimit = 0

// DefaultSessionName default sesion name for api
const DefaultSessionName = "mouthful-session"

// DefaultAuthorLengthLimit default author length limit
const DefaultAuthorLengthLimit = 50

// DefaultCleanupPeriod default cleanup period time
const DefaultCleanupPeriod = int64(86400)
