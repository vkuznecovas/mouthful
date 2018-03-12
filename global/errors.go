package global

import "errors"

// ErrThreadNotFound indicates that no thread was found
var ErrThreadNotFound = errors.New("Thread not found")

// ErrInternalServerError -> 500 basicly
var ErrInternalServerError = errors.New("Internal server error")

// ErrBadRequest -> 400 basicly
var ErrBadRequest = errors.New("Bad request")

// ErrUnauthorized -> 401 basicly
var ErrUnauthorized = errors.New("Unauthorized")

// ErrCommentNotFound indicates that the comment does not exist
var ErrCommentNotFound = errors.New("Comment not found")

// ErrWrongReplyTo indicates that the comments replyTo comment Id is invalid
var ErrWrongReplyTo = errors.New("Can't reply to this comment")
