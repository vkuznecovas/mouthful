package sqlite_test

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jmoiron/sqlx"
	"github.com/vkuznecovas/mouthful/db/sqlite"
	"github.com/vkuznecovas/mouthful/global"
)

func setupTestDb() *sqlite.Database {
	database := sqlite.Database{}
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	database.DB = db
	err = database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return &database
}

func TestCreateThread(t *testing.T) {
	database := setupTestDb()
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func TestCreateThreadUniqueViolation(t *testing.T) {
	database := setupTestDb()
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	uidNew, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, *uidNew))
}

func TestGetThread(t *testing.T) {
	database := setupTestDb()
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func TestGetThreadNotFound(t *testing.T) {
	database := setupTestDb()
	_, err := database.GetThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func TestCreateComment(t *testing.T) {
	now := time.Now().UTC()
	database := setupTestDb()
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
	database := setupTestDb()
	replyTo := global.GetUUID()
	_, err := database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestCreateCommentWithReply(t *testing.T) {
	database := setupTestDb()
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
}

func TestCreateCommentWrongReply(t *testing.T) {
	database := setupTestDb()
	uid1, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	uid2, err := database.CreateComment("body", "author", "/test", true, uid1)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid2)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestCreateCommentWrongThread(t *testing.T) {
	database := setupTestDb()
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/testasdasdasd", true, uid)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestGetCommentsByThreadNoThread(t *testing.T) {
	database := setupTestDb()
	_, err := database.GetCommentsByThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func TestGetCommentsByThreadEmptyThread(t *testing.T) {
	database := setupTestDb()
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	comments, err := database.GetCommentsByThread("/test")
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func TestGetCommentsByThread(t *testing.T) {
	database := setupTestDb()
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body1", "author1", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body2", "author2", "/test", false, nil)
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

func TestGetCommentNotFound(t *testing.T) {
	database := setupTestDb()
	_, err := database.GetComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestGetComment(t *testing.T) {
	database := setupTestDb()
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

func TestUpdateCommentNotFound(t *testing.T) {
	database := setupTestDb()
	err := database.UpdateComment(global.GetUUID(), "t", "t", false)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestUpdateComment(t *testing.T) {
	database := setupTestDb()
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

func TestDeleteCommentNotFound(t *testing.T) {
	database := setupTestDb()
	err := database.DeleteComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestDeleteComment(t *testing.T) {
	database := setupTestDb()
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	c, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.NotNil(t, c.DeletedAt)
}

func TestGetAllThreadsEmptyDatabase(t *testing.T) {
	database := setupTestDb()
	threads, err := database.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, threads, 0)
}

func TestGetAllThreads(t *testing.T) {
	database := setupTestDb()
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

func TestGetAllCommentsEmptyDatabase(t *testing.T) {
	database := setupTestDb()
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func TestGetAllComments(t *testing.T) {
	database := setupTestDb()
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

func TestSoftDelete(t *testing.T) {
	database := setupTestDb()
	author := "author"
	body := "body"
	path := "/test"
	uid, err := database.CreateComment(body, author, path, false, nil)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	c, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.NotNil(t, c.DeletedAt)
	err = database.RestoreDeletedComment(*uid)
	assert.Nil(t, err)
	c, err = database.GetComment(*uid)
	assert.Nil(t, err)
	assert.Nil(t, c.DeletedAt)
}

func TestGetAllCommentsGetsSoftDeletedComments(t *testing.T) {
	database := setupTestDb()
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

func TestDeleteCommentDeletesReplies(t *testing.T) {
	database := setupTestDb()
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
