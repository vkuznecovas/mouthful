package dynamodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/db"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
)

func setupDynamoTestDb() abstraction.Database {
	database := dynamodb.CreateTestDatabase()
	err := database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return database
}

func TestDynamoDb(t *testing.T) {
	testDb := setupDynamoTestDb()
	driver := testDb.GetUnderlyingStruct()
	driverCasted := driver.(*dynamodb.Database)
	for _, f := range db.TestFunctions {
		f.(func(*testing.T, abstraction.Database))(t, testDb)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}
