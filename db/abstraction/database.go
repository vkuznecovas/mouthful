// Package abstraction defines all the required interfaces to make the database layer pluggable.
package abstraction

import (
	"github.com/gofrs/uuid"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

// Database is a database instance for your selected DB
type Database interface {
	InitializeDatabase() error
	CreateThread(path string) (*uuid.UUID, error)
	GetThread(path string) (thread model.Thread, err error)
	CreateComment(body string, author string, path string, confirmed bool, replyTo *uuid.UUID) (*uuid.UUID, error)
	GetCommentsByThread(path string) ([]model.Comment, error)
	UpdateComment(id uuid.UUID, body, author string, confirmed bool) error
	DeleteComment(id uuid.UUID) error
	RestoreDeletedComment(id uuid.UUID) error
	GetComment(id uuid.UUID) (model.Comment, error)
	GetAllThreads() ([]model.Thread, error)
	GetAllComments() ([]model.Comment, error)
	GetDatabaseDialect() string
	GetUnderlyingStruct() interface{}
	CleanUpStaleData(target global.CleanupType, timeout int64) error
	HardDeleteComment(commentId uuid.UUID) error
	ImportData(pathToDump string) error
}
