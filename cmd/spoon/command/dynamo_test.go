package command_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/cmd/spoon/command"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"
)

func TestDynamoImport(t *testing.T) {
	sqlitePath := "./mouthful_test_db"
	sqliteCfg := model.Config{
		Database: model.Database{
			Dialect:  "sqlite3",
			Database: &sqlitePath,
		},
	}

	sqliteDb, err := sqlite.CreateDatabase(sqliteCfg.Database)
	assert.Nil(t, err)
	defer func() { os.Remove(sqlitePath) }()
	cid, err := sqliteDb.CreateComment("test", "test", "/test", true, nil)
	assert.Nil(t, err)
	cid, err = sqliteDb.CreateComment("test", "test", "/test", true, cid)
	assert.Nil(t, err)
	cid, err = sqliteDb.CreateComment("test", "test", "/testasdasd", true, nil)
	assert.Nil(t, err)

	cid, err = sqliteDb.CreateComment("test", "test", "/testasasdasddasd", true, nil)
	assert.Nil(t, err)

	err = sqliteDb.DeleteComment(*cid)
	assert.Nil(t, err)

	str := sqliteDb.GetUnderlyingStruct()
	strCasted := str.(*sqlxDriver.Database)
	err = strCasted.DB.Close()
	assert.Nil(t, err)
	wu := int64(1)
	test := "test"
	region := "eu-west-1"
	endpoint := "http://localhost:8000"
	dynamoDbConfig := model.Config{
		Database: model.Database{
			Dialect:                   "dynamodb",
			DynamoDBThreadReadUnits:   &wu,
			DynamoDBCommentReadUnits:  &wu,
			DynamoDBThreadWriteUnits:  &wu,
			DynamoDBCommentWriteUnits: &wu,
			DynamoDBIndexWriteUnits:   &wu,
			DynamoDBIndexReadUnits:    &wu,
			AwsRegion:                 &region,
			AwsAccessKeyID:            &test,
			AwsSecretAccessKey:        &test,
			DynamoDBEndpoint:          &endpoint,
		},
	}
	res, err := json.Marshal(dynamoDbConfig)
	assert.Nil(t, err)
	testFile := "./test-file"
	err = ioutil.WriteFile(testFile, res, 0644)
	assert.Nil(t, err)
	defer func() { os.Remove(testFile) }()

	err = command.DynamoCommandRun(sqlitePath, testFile)
	assert.Nil(t, err)
}
