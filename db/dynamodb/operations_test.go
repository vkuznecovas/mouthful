package dynamodb_test

import (
	"reflect"
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
	suiteReflected := reflect.ValueOf(db.TestSuite{})
	dbReflected := reflect.ValueOf(testDb)
	tReflected := reflect.ValueOf(t)
	for _, f := range db.TestFunctions {
		in := []reflect.Value{suiteReflected, tReflected, dbReflected}
		f.Call(in)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}

func TestDynamoDatawipe(t *testing.T) {
	testDb := setupDynamoTestDb()
	driver := testDb.GetUnderlyingStruct()
	driverCasted := driver.(*dynamodb.Database)
	body := "body"
	author := "author"
	path := "/"
	_, err := testDb.CreateComment(body, author, path, true, nil)
	assert.Nil(t, err)
	_, err = testDb.CreateComment(body, author, path, true, nil)
	assert.Nil(t, err)
	_, err = testDb.CreateThread("/t")
	assert.Nil(t, err)

	c, err := testDb.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, c, 2)
	th, err := testDb.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, th, 2)
	err = driverCasted.WipeOutData()
	assert.Nil(t, err)

	c, err = testDb.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, c, 0)

	th, err = testDb.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, th, 0)
}

func TestDynamoDialect(t *testing.T) {
	testDb := setupDynamoTestDb()
	assert.Equal(t, "dynamodb", testDb.GetDatabaseDialect())
}
