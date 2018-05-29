package global

import "errors"

// ErrThreadNotFound indicates that no thread was found
var ErrThreadNotFound = errors.New("Thread not found")

// ErrInternalServerError -> 500 basically
var ErrInternalServerError = errors.New("Internal server error")

// ErrBadRequest -> 400 basically
var ErrBadRequest = errors.New("Bad request")

// ErrUnauthorized -> 401 basically
var ErrUnauthorized = errors.New("Unauthorized")

// ErrCommentNotFound indicates that the comment does not exist
var ErrCommentNotFound = errors.New("Comment not found")

// ErrWrongReplyTo indicates that the comments replyTo comment Id is invalid
var ErrWrongReplyTo = errors.New("Can't reply to this comment")

// ErrCouldNotFindBundleFile indicates that we could not find the required bundle file
var ErrCouldNotFindBundleFile = errors.New("Can't find bundle file")

// ErrCouldNotOverrideBundlePath indicates that we could not find override the path in bundle
var ErrCouldNotOverrideBundlePath = errors.New("Can't override bundle file path")
