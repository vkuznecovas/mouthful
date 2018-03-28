package dynamodb_test

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
	"github.com/vkuznecovas/mouthful/global"
)

func setupTestDb() *abstraction.Database {
	tablePrefix := "something"
	database, err := dynamodb.CreateDatabase(model.Database{
		Dialect:     "dynamodb",
		TablePrefix: &tablePrefix,
	})
	if err != nil {
		panic(err)
	}
	wipeDB(database)
	err = database.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	return &database
}
func wipeDB(db abstraction.Database) {
	driver := db.GetUnderlyingStruct()
	driverCasted := driver.(*dynamodb.Database)
	_ = driverCasted.DB.Table(driverCasted.TablePrefix + global.DefaultDynamoDbThreadTableName).DeleteTable().Run()
	_ = driverCasted.DB.Table(driverCasted.TablePrefix + global.DefaultDynamoDbCommentTableName).DeleteTable().Run()
}

func TestDynamoDb(t *testing.T) {
	_ = setupTestDb()
}

func TestCreateThread(t *testing.T) {
	database := *setupTestDb()

	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func TestCreateThreadUniqueViolation(t *testing.T) {
	database := *setupTestDb()
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	uidNew, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, *uidNew))
}

func TestGetThread(t *testing.T) {
	database := *setupTestDb()
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func TestGetThreadNotFound(t *testing.T) {
	database := *setupTestDb()
	_, err := database.GetThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func TestCreateComment(t *testing.T) {
	now := time.Now().UTC()
	database := *setupTestDb()
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	comment, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.Nil(t, comment.DeletedAt)
	assert.True(t, uuid.Equal(*uid, comment.Id))
	assert.Equal(t, "body", comment.Body)
	assert.Equal(t, "author", comment.Author)
	assert.Equal(t, true, comment.Confirmed)
	assert.Equal(t, true, comment.CreatedAt.UTC().After(now))
	assert.Nil(t, comment.ReplyTo)

}

func TestCreateCommentNoReply(t *testing.T) {
	database := *setupTestDb()
	replyTo := global.GetUUID()
	_, err := database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestCreateCommentWithReply(t *testing.T) {
	database := *setupTestDb()
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
}

func TestCreateCommentWrongReply(t *testing.T) {
	database := *setupTestDb()
	uid1, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	uid2, err := database.CreateComment("body", "author", "/test", true, uid1)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid2)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestCreateCommentWrongThread(t *testing.T) {
	database := *setupTestDb()
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/testasdasdasd", true, uid)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}
