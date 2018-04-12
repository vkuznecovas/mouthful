package db_test

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/mysql"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"

	"github.com/vkuznecovas/mouthful/global"
)

var testFunctions = [...]interface{}{CreateThread,
	CreateThreadUniqueViolation,
	GetThread,
	GetThreadNotFound,
	CreateComment,
	CreateCommentNoReply,
	CreateCommentWithReply,
	CreateCommentWrongReply,
	CreateCommentWrongThread,
	GetCommentNotFound,
	GetComment,
	GetCommentsByThreadNoThread,
	GetCommentsByThread,
	UpdateCommentNotFound,
	UpdateComment,
	DeleteCommentNotFound,
	DeleteComment,
	GetAllThreadsEmptyDatabase,
	GetAllThreads,
	GetAllCommentsEmptyDatabase,
	GetAllComments,
	SoftDelete,
	GetAllCommentsGetsSoftDeletedComments,
	DeleteCommentDeletesReplies,
	CreateCommentReplyToAReply,
}

func setupDynamoTestDb() abstraction.Database {
	database := dynamodb.CreateTestDatabase()
	err := database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return database
}

func setupSqliteTestDb() abstraction.Database {
	database := sqlite.CreateTestDatabase()
	return database
}

func setupMysqlTestDb() abstraction.Database {
	database := mysql.CreateTestDatabase()
	return database
}

func TestDynamoDb(t *testing.T) {
	db := setupDynamoTestDb()
	driver := db.GetUnderlyingStruct()
	driverCasted := driver.(*dynamodb.Database)
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, db)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}

func TestSqliteDb(t *testing.T) {
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, setupSqliteTestDb())
	}
}
func TestMysqlDb(t *testing.T) {
	db := mysql.CreateTestDatabase()
	driver := db.GetUnderlyingStruct()
	driverCasted := driver.(*sqlxDriver.Database)
	// clean out before start
	driverCasted.WipeOutData()
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, db)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}

func CreateThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func CreateThreadUniqueViolation(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	uidNew, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, *uidNew))
}

func GetThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func GetThreadNotFound(t *testing.T, database abstraction.Database) {
	_, err := database.GetThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func CreateComment(t *testing.T, database abstraction.Database) {
	now := time.Now().UTC()
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

func CreateCommentNoReply(t *testing.T, database abstraction.Database) {
	replyTo := global.GetUUID()
	_, err := database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func CreateCommentWithReply(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
}

func CreateCommentWrongReply(t *testing.T, database abstraction.Database) {
	_, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	uid2, err := database.CreateComment("body", "author", "/test1", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid2)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func CreateCommentReplyToAReply(t *testing.T, database abstraction.Database) {
	uid1, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	uid2, err := database.CreateComment("body", "author", "/test", true, uid1)
	assert.Nil(t, err)
	uid3, err := database.CreateComment("body", "author", "/test", true, uid2)
	assert.Nil(t, err)
	comment, err := database.GetComment(*uid3)
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*comment.ReplyTo, *uid1))
}

func CreateCommentWrongThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/testasdasdasd", true, uid)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func GetCommentNotFound(t *testing.T, database abstraction.Database) {
	_, err := database.GetComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func GetComment(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	comment, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, comment.Id))
	assert.Equal(t, "body", comment.Body)
	assert.Equal(t, true, comment.Confirmed)
	assert.Equal(t, "author", comment.Author)
	assert.Nil(t, comment.ReplyTo)
}

func GetCommentsByThreadNoThread(t *testing.T, database abstraction.Database) {
	_, err := database.GetCommentsByThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func GetCommentsByThreadEmptyThread(t *testing.T, database abstraction.Database) {
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	comments, err := database.GetCommentsByThread("/test")
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func GetCommentsByThread(t *testing.T, database abstraction.Database) {
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body1", "author1", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body2", "author2", "/test", false, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body2", "author2", "/test2", false, nil)
	assert.Nil(t, err)
	comments, err := database.GetCommentsByThread("/test")
	assert.Nil(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, "body", comments[0].Body)
	assert.Equal(t, "body1", comments[1].Body)
	assert.Equal(t, "author", comments[0].Author)
	assert.Equal(t, "author1", comments[1].Author)
	assert.Nil(t, comments[0].ReplyTo)
	assert.Nil(t, comments[1].ReplyTo)
	assert.Nil(t, comments[0].DeletedAt)
	assert.Nil(t, comments[1].DeletedAt)

	assert.Equal(t, true, comments[0].Confirmed)
	assert.Equal(t, true, comments[1].Confirmed)
}
func UpdateCommentNotFound(t *testing.T, database abstraction.Database) {
	err := database.UpdateComment(global.GetUUID(), "t", "t", false)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func UpdateComment(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.UpdateComment(*uid, "t", "t", false)
	assert.Nil(t, err)
	comment, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, comment.Id))
	assert.Equal(t, "t", comment.Body)
	assert.Equal(t, false, comment.Confirmed)
	assert.Equal(t, "t", comment.Author)
	assert.Nil(t, comment.ReplyTo)
}

func DeleteCommentNotFound(t *testing.T, database abstraction.Database) {
	err := database.DeleteComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func DeleteComment(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	c, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.NotNil(t, c.DeletedAt)
}

func GetAllThreadsEmptyDatabase(t *testing.T, database abstraction.Database) {
	threads, err := database.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, threads, 0)
}

func GetAllThreads(t *testing.T, database abstraction.Database) {
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	_, err = database.CreateThread("/test1")
	assert.Nil(t, err)
	threads, err := database.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, threads, 2)
	assert.Equal(t, "/test", threads[0].Path)
	assert.Equal(t, "/test1", threads[1].Path)
}

func GetAllCommentsEmptyDatabase(t *testing.T, database abstraction.Database) {
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func GetAllComments(t *testing.T, database abstraction.Database) {
	author := "author"
	body := "body"
	path := "/test"
	_, err := database.CreateComment(body, author, path, false, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment(body, author, path, true, nil)
	assert.Nil(t, err)
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, body, comments[0].Body)
	assert.Equal(t, body, comments[1].Body)
	assert.Equal(t, false, comments[0].Confirmed)
	assert.Equal(t, true, comments[1].Confirmed)
	assert.Equal(t, author, comments[0].Author)
	assert.Equal(t, author, comments[1].Author)
	assert.Nil(t, comments[0].ReplyTo)
	assert.Nil(t, comments[1].ReplyTo)
}

func SoftDelete(t *testing.T, database abstraction.Database) {
	author := "author"
	body := "body"
	path := "/test"
	uid, err := database.CreateComment(body, author, path, false, nil)
	assert.Nil(t, err)
	c, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.Nil(t, c.DeletedAt)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	c, err = database.GetComment(*uid)
	assert.Nil(t, err)
	assert.NotNil(t, c.DeletedAt)
	err = database.RestoreDeletedComment(*uid)
	assert.Nil(t, err)
	c, err = database.GetComment(*uid)
	assert.Nil(t, err)
	assert.Nil(t, c.DeletedAt)
}

func GetAllCommentsGetsSoftDeletedComments(t *testing.T, database abstraction.Database) {
	author := "author"
	body := "body"
	path := "/test"
	uid, err := database.CreateComment(body, author, path, false, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment(body, author, path, true, nil)
	assert.Nil(t, err)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, body, comments[0].Body)
	assert.Equal(t, body, comments[1].Body)
	assert.Equal(t, false, comments[0].Confirmed)
	assert.Equal(t, true, comments[1].Confirmed)
	assert.Equal(t, author, comments[0].Author)
	assert.Equal(t, author, comments[1].Author)
	assert.NotNil(t, comments[0].DeletedAt)
	assert.Nil(t, comments[1].DeletedAt)
	assert.Nil(t, comments[0].ReplyTo)
	assert.Nil(t, comments[1].ReplyTo)
}

func DeleteCommentDeletesReplies(t *testing.T, database abstraction.Database) {
	author := "author"
	body := "body"
	path := "/test"
	uid, err := database.CreateComment(body, author, path, false, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment(body, author, path, true, uid)
	assert.Nil(t, err)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 2)
	assert.NotNil(t, comments[0].DeletedAt)
	assert.NotNil(t, comments[1].DeletedAt)
}
