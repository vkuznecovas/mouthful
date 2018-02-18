package sqlite_test

import (
	"testing"
	"time"

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
	err := database.CreateThread("/test")
	assert.Nil(t, err)
	rows, err := database.DB.Query("select * from thread limit 1")
	assert.Nil(t, err)
	assert.NotNil(t, rows)
	for rows.Next() {
		id := 0
		path := ""
		err = rows.Scan(&id, &path)
		assert.Nil(t, err)
		assert.Equal(t, 1, id)
		assert.Equal(t, "/test", path)
	}
}

func TestGetThread(t *testing.T) {
	database := setupTestDb()
	err := database.CreateThread("/test")
	assert.Nil(t, err)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.Equal(t, 1, thread.Id)
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
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	rows, err := database.DB.Query("select * from comment limit 1")
	assert.Nil(t, err)
	assert.NotNil(t, rows)
	for rows.Next() {
		id := 0
		threadId := 0
		body := ""
		author := ""
		confirmed := false
		replyTo := new(*int)
		createdAt := time.Now()
		err = rows.Scan(&id, &threadId, &body, &author, &confirmed, &createdAt, &replyTo)
		assert.Nil(t, err)
		assert.Equal(t, 1, id)
		assert.Equal(t, 1, threadId)
		assert.Equal(t, "body", body)
		assert.Equal(t, "author", author)
		assert.Equal(t, true, confirmed)
		assert.Equal(t, true, createdAt.UTC().After(now))
		assert.Nil(t, replyTo)
	}
}

func TestCreateCommentNoReply(t *testing.T) {
	database := setupTestDb()
	replyTo := 1
	err := database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestCreateCommentWithReply(t *testing.T) {
	database := setupTestDb()
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	replyTo := 1
	err = database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.Nil(t, err)
}

func TestCreateCommentWrongReply(t *testing.T) {
	database := setupTestDb()
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	replyTo := 1
	err = database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.Nil(t, err)
	replyTo = 2
	err = database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func TestCreateCommentWrongThread(t *testing.T) {
	database := setupTestDb()
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	replyTo := 1
	err = database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.Nil(t, err)
	err = database.CreateComment("body", "author", "/testasdasdasd", true, &replyTo)
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
	err := database.CreateThread("/test")
	assert.Nil(t, err)
	comments, err := database.GetCommentsByThread("/test")
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func TestGetCommentsByThread(t *testing.T) {
	database := setupTestDb()
	err := database.CreateThread("/test")
	assert.Nil(t, err)
	err = database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.CreateComment("body1", "author1", "/test", true, nil)
	assert.Nil(t, err)
	err = database.CreateComment("body2", "author2", "/test", false, nil)
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

	assert.Equal(t, true, comments[0].Confirmed)
	assert.Equal(t, true, comments[1].Confirmed)
}

func TestGetCommentNotFound(t *testing.T) {
	database := setupTestDb()
	_, err := database.GetComment(1)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestGetComment(t *testing.T) {
	database := setupTestDb()
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	comment, err := database.GetComment(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, comment.Id)
	assert.Equal(t, "body", comment.Body)
	assert.Equal(t, true, comment.Confirmed)
	assert.Equal(t, "author", comment.Author)
	assert.Nil(t, comment.ReplyTo)
}

func TestUpdateCommentNotFound(t *testing.T) {
	database := setupTestDb()
	err := database.UpdateComment(1, "t", "t", false)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestUpdateComment(t *testing.T) {
	database := setupTestDb()
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.UpdateComment(1, "t", "t", false)
	assert.Nil(t, err)
	comment, err := database.GetComment(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, comment.Id)
	assert.Equal(t, "t", comment.Body)
	assert.Equal(t, false, comment.Confirmed)
	assert.Equal(t, "t", comment.Author)
	assert.Nil(t, comment.ReplyTo)
}

func TestDeleteCommentNotFound(t *testing.T) {
	database := setupTestDb()
	err := database.DeleteComment(1)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestDeleteComment(t *testing.T) {
	database := setupTestDb()
	err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.DeleteComment(1)
	assert.Nil(t, err)
	_, err = database.GetComment(1)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func TestGetAllThreadsEmptyDatabase(t *testing.T) {
	database := setupTestDb()
	threads, err := database.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, threads, 0)
}

func TestGetAllThreads(t *testing.T) {
	database := setupTestDb()
	err := database.CreateThread("/test")
	assert.Nil(t, err)
	err = database.CreateThread("/test1")
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
	err := database.CreateComment(body, author, path, false, nil)
	assert.Nil(t, err)
	err = database.CreateComment(body, author, path, true, nil)
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
