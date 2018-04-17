// Package global provides constants, defaults and helper methods to mouthful.
package global

// StaticPath is the default location for our static files
const StaticPath = "./static"

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
