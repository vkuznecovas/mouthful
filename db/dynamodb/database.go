package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	dynamoModel "github.com/vkuznecovas/mouthful/db/dynamodb/model"
	"github.com/vkuznecovas/mouthful/global"
)

// TODO: tests

// Database is a database instance for sqlite
type Database struct {
	DB          *dynamo.DB
	Config      model.Database
	TablePrefix string
	IsTest      bool
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
	prefix := global.GetUUID().String() + "_"
	database := &Database{
		DB: db,
		Config: model.Database{
			TablePrefix: &prefix,
		},
		IsTest: true,
	}
	err := database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return database
}

// WipeOutData deletes all the threads and comments in the database if the database is a test one
func (d *Database) WipeOutData() {
	if !d.IsTest {
		return
	}
	var threads []dynamoModel.Thread
	var comments []dynamoModel.Comment
	err := d.DB.Table(d.TablePrefix + global.DefaultDynamoDbThreadTableName).Scan().All(&threads)
	if err != nil {
		panic(err)
	}
	for _, v := range threads {
		err := d.DB.Table(d.TablePrefix+global.DefaultDynamoDbThreadTableName).Delete("Path", v.Path).Run()
		if err != nil {
			panic(err)
		}
	}
	err = d.DB.Table(d.TablePrefix + global.DefaultDynamoDbCommentTableName).Scan().All(&comments)
	if err != nil {
		panic(err)
	}
	for _, v := range comments {
		err := d.DB.Table(d.TablePrefix+global.DefaultDynamoDbCommentTableName).Delete("ID", v.Id).Run()
		if err != nil {
			panic(err)
		}
	}
}

// DeleteTables deletes the thread and comment tables in the database if the database is a test one
func (d *Database) DeleteTables() {
	if !d.IsTest {
		return
	}
	err := d.DB.Table(d.TablePrefix + global.DefaultDynamoDbThreadTableName).DeleteTable().Run()
	if err != nil {
		panic(err)
	}
	err = d.DB.Table(d.TablePrefix + global.DefaultDynamoDbCommentTableName).DeleteTable().Run()
	if err != nil {
		panic(err)
	}
}
