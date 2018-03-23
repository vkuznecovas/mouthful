package dynamodb

import (
	"github.com/satori/go.uuid"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

// InitializeDatabase runs the queries for an initial database seed
func (db *Database) InitializeDatabase() error {
	tables := [...]string{"thread", "comment"}
	for i := range tables {
		tables[i] = db.TablePrefix + tables[i]
	}

	dynamoTables, err := db.DB.ListTables().All()
	if err != nil {
		return err
	}
	for _, t := range tables {
		found := false

		for _, v := range dynamoTables {

			if v == t {
				found = true
			}

		}
		if !found {
			db.DB.CreateTable(t)
		}
	}

	return nil
}

// CreateThread takes the thread path and creates it in the database
func (db *Database) CreateThread(path string) (*uuid.UUID, error) {
	uid := global.GetUUID()
	return &uid, nil
}

// GetThread takes the thread path and fetches it from the database
func (db *Database) GetThread(path string) (thread model.Thread, err error) {

	return thread, err
}

// CreateComment takes in a body, author, and path and creates a comment for the given thread. If thread does not exist, it creates one
func (db *Database) CreateComment(body string, author string, path string, confirmed bool, replyTo *uuid.UUID) (*uuid.UUID, error) {

	uid := global.GetUUID()
	return &uid, nil
}

// GetCommentsByThread gets all the comments by thread path
func (db *Database) GetCommentsByThread(path string) (comments []model.Comment, err error) {

	return comments, nil
}

// GetComment gets comment by id
func (db *Database) GetComment(id uuid.UUID) (comment model.Comment, err error) {

	return comment, nil
}

// UpdateComment updatesComment comment by id
func (db *Database) UpdateComment(id uuid.UUID, body, author string, confirmed bool) error {

	return nil
}

// DeleteComment soft-deletes the comment by id and all the replies to it
func (db *Database) DeleteComment(id uuid.UUID) error {

	return nil
}

// RestoreDeletedComment restores the soft-deleted comment
func (db *Database) RestoreDeletedComment(id uuid.UUID) error {

	return nil
}

// GetAllThreads gets all the threads found in the database
func (db *Database) GetAllThreads() (threads []model.Thread, err error) {
	return threads, err
}

// GetAllComments gets all the comments found in the database
func (db *Database) GetAllComments() (comments []model.Comment, err error) {
	return comments, err
}

// GetDatabaseDialect returns the current database dialect
func (db *Database) GetDatabaseDialect() string {
	return "dynamodb"
}
