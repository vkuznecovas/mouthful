package dynamodb_test

import (
	"github.com/vkuznecovas/mouthful/db/dynamodb"
)

func setupTestDb() *dynamodb.Database {
	database := dynamodb.CreateDatabase()
	err = database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return &database
}
