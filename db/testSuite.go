package db

import (
	"reflect"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/db/abstraction"

	"github.com/vkuznecovas/mouthful/global"
)

// TestFunctions is a list of database test functions used to test the drivers.
// If you want to add a test method, just follow this signature
// func (ts TestSuite) SomeName(t *testing.T, database abstraction.Database)
// It will then get automagically added to the test functions by init method.
var TestFunctions []reflect.Value

func init() {
	testSuiteType := reflect.TypeOf(TestSuite{})
	for i := 0; i < testSuiteType.NumMethod(); i++ {
		method := testSuiteType.Method(i)
		TestFunctions = append(TestFunctions, method.Func)
	}
}

// TestSuite contains all the functions that one needs to test database operations against.
type TestSuite struct {
}

// CreateThread checks if a thread is correctly created
func (ts TestSuite) CreateThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

// CreateThreadUniqueViolation checks if duplicate thread creation throws no errors
func (ts TestSuite) CreateThreadUniqueViolation(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	uidNew, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, *uidNew))
}

// GetThread checks if a created thread is gotten alright
func (ts TestSuite) GetThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateThread("/test")
	assert.Nil(t, err)
	assert.NotNil(t, uid)
	thread, err := database.GetThread("/test")
	assert.Nil(t, err)
	assert.True(t, uuid.Equal(*uid, thread.Id))
	assert.Equal(t, "/test", thread.Path)
}

// GetThreadNotFound asserts that we correctly get a response saying we're not finding the thread
func (ts TestSuite) GetThreadNotFound(t *testing.T, database abstraction.Database) {
	_, err := database.GetThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

// CreateComment checks if we create the comment alright
func (ts TestSuite) CreateComment(t *testing.T, database abstraction.Database) {
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

// CreateCommentNoReply checks if we return an error upon replying to a non existant reply to
func (ts TestSuite) CreateCommentNoReply(t *testing.T, database abstraction.Database) {
	replyTo := global.GetUUID()
	_, err := database.CreateComment("body", "author", "/test", true, &replyTo)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

// CreateCommentWithReply checks if a comment with a proper reply to value gets created correctly
func (ts TestSuite) CreateCommentWithReply(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
}

// CreateCommentWrongReply asserts that upon providing a bad reply to value we turn a ErrWrongReplyTo
func (ts TestSuite) CreateCommentWrongReply(t *testing.T, database abstraction.Database) {
	_, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	uid2, err := database.CreateComment("body", "author", "/test1", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid2)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

// CreateCommentReplyToAReply asserts that we only allow 2 levels of nesting by making sure that a reply to a reply will instead point to it's parent
func (ts TestSuite) CreateCommentReplyToAReply(t *testing.T, database abstraction.Database) {
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

// CreateCommentWrongThread asserts that we return an error upon trying to reply to a comment from another thread
func (ts TestSuite) CreateCommentWrongThread(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/test", true, uid)
	assert.Nil(t, err)
	_, err = database.CreateComment("body", "author", "/testasdasdasd", true, uid)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrWrongReplyTo, err)
}

// GetCommentNotFound asserts that we return ErrCommentNotFound if a comment is not found on GetComment
func (ts TestSuite) GetCommentNotFound(t *testing.T, database abstraction.Database) {
	_, err := database.GetComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

// GetComment checks if the getComment operation actually gets the comment
func (ts TestSuite) GetComment(t *testing.T, database abstraction.Database) {
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

// GetCommentsByThreadNoThread asserts that we return ErrThreadNotFound if no thread is found
func (ts TestSuite) GetCommentsByThreadNoThread(t *testing.T, database abstraction.Database) {
	_, err := database.GetCommentsByThread("/test")
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrThreadNotFound, err)
}

// GetCommentsByThreadEmptyThread asserts that we return an empty array if the thread has no comments
func (ts TestSuite) GetCommentsByThreadEmptyThread(t *testing.T, database abstraction.Database) {
	_, err := database.CreateThread("/test")
	assert.Nil(t, err)
	comments, err := database.GetCommentsByThread("/test")
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

// GetCommentsByThread asserts that we get correct comments for a specific thread, aka only confirmed ones.
func (ts TestSuite) GetCommentsByThread(t *testing.T, database abstraction.Database) {
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

// UpdateCommentNotFound asserts that we return ErrCommentNotFound upon updating a non existant comment
func (ts TestSuite) UpdateCommentNotFound(t *testing.T, database abstraction.Database) {
	err := database.UpdateComment(global.GetUUID(), "t", "t", false)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

// UpdateComment checks if we update the comment alright
func (ts TestSuite) UpdateComment(t *testing.T, database abstraction.Database) {
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

// DeleteCommentNotFound asserts if ErrCommentNotFound is return upon deletion of a non existant comment
func (ts TestSuite) DeleteCommentNotFound(t *testing.T, database abstraction.Database) {
	err := database.DeleteComment(global.GetUUID())
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCommentNotFound, err)
}

// DeleteComment asserts that we actually soft delete a comment
func (ts TestSuite) DeleteComment(t *testing.T, database abstraction.Database) {
	uid, err := database.CreateComment("body", "author", "/test", true, nil)
	assert.Nil(t, err)
	err = database.DeleteComment(*uid)
	assert.Nil(t, err)
	c, err := database.GetComment(*uid)
	assert.Nil(t, err)
	assert.NotNil(t, c.DeletedAt)
}

// GetAllThreadsEmptyDatabase asserts that no threads are returned if none exist
func (ts TestSuite) GetAllThreadsEmptyDatabase(t *testing.T, database abstraction.Database) {
	threads, err := database.GetAllThreads()
	assert.Nil(t, err)
	assert.Len(t, threads, 0)
}

// GetAllThreads asserts that we return all the threads correctly
func (ts TestSuite) GetAllThreads(t *testing.T, database abstraction.Database) {
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

// GetAllCommentsEmptyDatabase asserts that we return an empty dataset
func (ts TestSuite) GetAllCommentsEmptyDatabase(t *testing.T, database abstraction.Database) {
	comments, err := database.GetAllComments()
	assert.Nil(t, err)
	assert.Len(t, comments, 0)
}

// GetAllComments asserts that all comments are gotten correctly.
func (ts TestSuite) GetAllComments(t *testing.T, database abstraction.Database) {
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

// SoftDelete checks if comments are soft deleted
func (ts TestSuite) SoftDelete(t *testing.T, database abstraction.Database) {
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

// GetAllCommentsGetsSoftDeletedComments checks if get all comments returns all the comments, even the deleted ones.
func (ts TestSuite) GetAllCommentsGetsSoftDeletedComments(t *testing.T, database abstraction.Database) {
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

// DeleteCommentDeletesReplies asserts that deletes are cascaded
func (ts TestSuite) DeleteCommentDeletesReplies(t *testing.T, database abstraction.Database) {
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
