package abstraction

import (
	"github.com/vkuznecovas/mouthful/db/model"
)

// Db is a database instance for your selected DB
type Database interface {
	InitializeDatabase() error
	CreateThread(path string) error
	CreateComment(body string, author string, path string, confirmed bool, replyTo *int) error
	GetCommentsByThread(path string) ([]model.Comment, error)
	UpdateComment(id int, body, author string, confirmed bool) error
	DeleteComment(id int) error
	RestoreDeletedComment(id int) error
	GetComment(id int) (model.Comment, error)
	GetAllThreads() ([]model.Thread, error)
	GetAllComments() ([]model.Comment, error)
	GetDatabaseDialect() string
}
