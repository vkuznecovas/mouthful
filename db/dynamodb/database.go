package dynamodb

import (
	"errors"
	"os"

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

// ValidateConfig validates the config for sqlite
func ValidateConfig(config model.Database) error {
	err := ""
	// TODO: write recommended units here
	if config.DynamoDBCommentReadUnits == nil {
		err += "Please specify the read units for dynamoDb Comment table by adjusting the config value DynamoDBCommentReadUnits\n"
	}
	if config.DynamoDBThreadReadUnits == nil {
		err += "Please specify the read units for dynamoDb Thread table by adjusting the config value DynamoDBThreadReadUnits\n"
	}
	if config.DynamoDBThreadWriteUnits == nil {
		err += "Please specify the write units for dynamoDb Thread table by adjusting the config value DynamoDBThreadWriteUnits\n"
	}
	if config.DynamoDBCommentWriteUnits == nil {
		err += "Please specify the write units for dynamoDb Comment table by adjusting the config value DynamoDBCommentWriteUnits\n"
	}
	if config.DynamoDBIndexReadUnits == nil {
		err += "Please specify the write units for dynamoDb Thread table by adjusting the config value DynamoDBThreadWriteUnits\n"
	}
	if config.DynamoDBIndexWriteUnits == nil {
		err += "Please specify the write units for dynamoDb Comment table by adjusting the config value DynamoDBIndexWriteUnits\n"
	}
	if config.AwsRegion == nil {
		err += "Please specify the AWS region your dynamoDb lives in\n"
	}
	if config.AwsAccessKeyID == nil {
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
			err += "Please specify the AWS_ACCESS_KEY_ID either by setting AWS_ACCESS_KEY_ID environment variable or setting AwsAccessKeyID in config\n"
		} else {
			value := os.Getenv("AWS_ACCESS_KEY_ID")
			config.AwsAccessKeyID = &value
		}
	} else {
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
			os.Setenv("AWS_ACCESS_KEY_ID", *config.AwsAccessKeyID)
		}
	}
	if config.AwsSecretAccessKey == nil {
		if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
			err += "Please specify the AWS_SECRET_ACCESS_KEY either by setting AWS_SECRET_ACCESS_KEY environment variable or setting AwsSecretAccessKey in config\n"
		}
	} else {
		if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
			os.Setenv("AWS_SECRET_ACCESS_KEY", *config.AwsSecretAccessKey)
		}
	}
	if err != "" {
		return errors.New(err)
	}
	return nil
}

// CreateDatabase creates a database instance from the given config
func CreateDatabase(databaseConfig model.Database) (abstraction.Database, error) {
	err := ValidateConfig(databaseConfig)
	if err != nil {
		return nil, err
	}
	cfg := &aws.Config{Region: aws.String(*databaseConfig.AwsRegion)}
	if databaseConfig.DynamoDBEndpoint != nil {
		cfg.Endpoint = aws.String(*databaseConfig.DynamoDBEndpoint)
	}

	db := dynamo.New(session.New(), cfg)
	prefix := ""
	if databaseConfig.TablePrefix != nil {
		prefix = *databaseConfig.TablePrefix
	}

	newDb := &Database{
		DB:          db,
		Config:      databaseConfig,
		TablePrefix: prefix,
	}
	err = newDb.InitializeDatabase()
	if err != nil {
		return nil, err
	}
	return newDb, nil
}

// CreateTestDatabase creates a database instance for testing locally.
// It creates tables with UUID prefix, so should be safe to use even if tests are run in parallel.
func CreateTestDatabase() abstraction.Database {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("eu-west-1"), Endpoint: aws.String("http://localhost:8000")})
	prefix := global.GetUUID().String() + "_"
	units := int64(1)
	database := &Database{
		DB: db,
		Config: model.Database{
			TablePrefix:               &prefix,
			DynamoDBCommentReadUnits:  &units,
			DynamoDBCommentWriteUnits: &units,
			DynamoDBThreadReadUnits:   &units,
			DynamoDBThreadWriteUnits:  &units,
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
func (d *Database) WipeOutData() error {
	if !d.IsTest {
		return nil
	}
	var threads []dynamoModel.Thread
	var comments []dynamoModel.Comment
	err := d.DB.Table(d.TablePrefix + global.DefaultDynamoDbThreadTableName).Scan().All(&threads)
	if err != nil {
		return err
	}
	for _, v := range threads {
		err := d.DB.Table(d.TablePrefix+global.DefaultDynamoDbThreadTableName).Delete("Path", v.Path).Run()
		if err != nil {
			return err
		}
	}
	err = d.DB.Table(d.TablePrefix + global.DefaultDynamoDbCommentTableName).Scan().All(&comments)
	if err != nil {
		return err
	}
	for _, v := range comments {
		err := d.DB.Table(d.TablePrefix+global.DefaultDynamoDbCommentTableName).Delete("ID", v.Id).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteTables deletes the thread and comment tables in the database if the database is a test one
func (d *Database) DeleteTables() error {
	if !d.IsTest {
		return nil
	}
	err := d.DB.Table(d.TablePrefix + global.DefaultDynamoDbThreadTableName).DeleteTable().Run()
	if err != nil {
		return err
	}
	err = d.DB.Table(d.TablePrefix + global.DefaultDynamoDbCommentTableName).DeleteTable().Run()
	if err != nil {
		return err
	}
	return nil
}
