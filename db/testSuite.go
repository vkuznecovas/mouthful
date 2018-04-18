package db

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/db/abstraction"

	"github.com/vkuznecovas/mouthful/global"
)

// TestFunctions is a list of database test functions used to test the drivers
var TestFunctions = [...]interface{}{createThread,
	createThreadUniqueViolation,
	getThread,
	getThreadNotFound,
	createComment,
	createCommentNoReply,
	createCommentWithReply,
	createCommentWrongReply,
	createCommentWrongThread,
	getCommentNotFound,
	getComment,
	getCommentsByThreadNoThread,
	getCommentsByThread,
	updateCommentNotFound,
	updateComment,
	deleteCommentNotFound,
	deleteComment,
	getAllThreadsEmptyDatabase,
	getAllThreads,
	getAllCommentsEmptyDatabase,
	getAllComments,
	softDelete,
	getAllCommentsGetsSoftDeletedComments,
	deleteCommentDeletesReplies,
	createCommentReplyToAReply,
}

func createThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func createThreadUniqueViolation(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	uidNew, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, *uidNew))
}

func getThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

func getThreadNotFound(t *testing.T, database abstraction.Database) {
	_, err := database.GetThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func createComment(t *testing.T, database abstraction.Database) {
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

func createCommentNoReply(t *testing.T, database abstraction.Database) {
	replyTo := global.GetUUID()
	_, err := database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func createCommentWithReply(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
}

func createCommentWrongReply(t *testing.T, database abstraction.Database) {
	_, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	uid2, err := database.CreateComment("body", "author", "/test1", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid2)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func createCommentReplyToAReply(t *testing.T, database abstraction.Database) {
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

func createCommentWrongThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/testasdasdasd", true, uid)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

func getCommentNotFound(t *testing.T, database abstraction.Database) {
	_, err := database.GetComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func getComment(t *testing.T, database abstraction.Database) {
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

func getCommentsByThreadNoThread(t *testing.T, database abstraction.Database) {
	_, err := database.GetCommentsByThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

func getCommentsByThreadEmptyThread(t *testing.T, database abstraction.Database) {
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	comments, err := database.GetCommentsByThread("/test")
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func getCommentsByThread(t *testing.T, database abstraction.Database) {
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
func updateCommentNotFound(t *testing.T, database abstraction.Database) {
	err := database.UpdateComment(global.GetUUID(), "t", "t", false)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func updateComment(t *testing.T, database abstraction.Database) {
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

func deleteCommentNotFound(t *testing.T, database abstraction.Database) {
	err := database.DeleteComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

func deleteComment(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	c, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.NotNil(t, c.DeletedAt)
}

func getAllThreadsEmptyDatabase(t *testing.T, database abstraction.Database) {
	threads, err := database.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, threads, 0)
}

func getAllThreads(t *testing.T, database abstraction.Database) {
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

func getAllCommentsEmptyDatabase(t *testing.T, database abstraction.Database) {
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

func getAllComments(t *testing.T, database abstraction.Database) {
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

func softDelete(t *testing.T, database abstraction.Database) {
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

func getAllCommentsGetsSoftDeletedComments(t *testing.T, database abstraction.Database) {
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

func deleteCommentDeletesReplies(t *testing.T, database abstraction.Database) {
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
