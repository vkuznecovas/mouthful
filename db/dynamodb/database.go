package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/global"
)

// TODO: tests

// Database is a database instance for sqlite
type Database struct {
	DB          *dynamo.DB
	Config      model.Database
	TablePrefix string
}

// CreateDatabase creates a database instance from the given config
func CreateDatabase(databaseConfig model.Database) (abstraction.Database, error) {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("eu-west-1"), Endpoint: aws.String("http://localhost:8000")})
	prefix := ""
	if databaseConfig.TablePrefix != nil {
		prefix = *databaseConfig.TablePrefix
	}
	return &Database{
		DB:          db,
		Config:      databaseConfig,
		TablePrefix: prefix,
	}, nil
}

// CreateTestDatabase creates a database instance for testing locally.
// It creates tables with UUID prefix, so should be safe to use even if tests are run in parallel.
func CreateTestDatabase() abstraction.Database {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("eu-west-1"), Endpoint: aws.String("http://localhost:8000")})
	prefix := global.GetUUID().String()
	database := &Database{
		DB: db,
		Config: model.Database{
			TablePrefix: &prefix,
		},
	}
	err := database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return database
}
