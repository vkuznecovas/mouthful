package dynamodb

import (
	"log"
	"strings"
	"time"

	"github.com/guregu/dynamo"

	"github.com/satori/go.uuid"
	dynamoModel "github.com/vkuznecovas/mouthful/db/dynamodb/model"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

// InitializeDatabase runs the queries for an initial database seed
func (db *Database) InitializeDatabase() error {
	tables := [...]string{global.DefaultDynamoDbThreadTableName, global.DefaultDynamoDbCommentTableName}
	tableModelMap := map[string]interface{}{
		global.DefaultDynamoDbThreadTableName:  dynamoModel.Thread{},
		global.DefaultDynamoDbCommentTableName: model.Comment{},
	}

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
			log.Printf("Creating table %v\n", t)
			noPrefix := strings.Replace(t, db.TablePrefix, "", 1)
			err := db.DB.CreateTable(t, tableModelMap[noPrefix]).Provision(4, 2).Run()
			if err != nil {
				return err
			}
		}
	}

	log.Printf("Tables created, waiting for them to be ready. Timeout - 1 minute\n")
	// TODO check if this actually works in aws, it might not.
	for i := 0; i < 60; i++ {
		dt, err := db.DB.ListTables().All()
		if err != nil {
			return err
		}
		matches := 0
		for _, v := range dt {
			for _, t := range tables {
				if t == v {
					matches++
				}
			}
		}
		if matches == len(tables) {
			log.Printf("Tables created, continuing...\n")
			break
		} else {
			log.Printf("Waiting for tables...\n")
			time.Sleep(time.Second)
		}
	}

	return nil
}

// CreateThread takes the thread path and creates it in the database
func (db *Database) CreateThread(path string) (*uuid.UUID, error) {
	thread, err := db.GetThread(path)
	if err != nil {
		if err == global.ErrThreadNotFound {
			uid := global.GetUUID()
			err := db.DB.Table(db.TablePrefix + global.DefaultDynamoDbThreadTableName).Put(dynamoModel.Thread{
				Id:        uid,
				Path:      path,
				CreatedAt: time.Now(),
				Comments:  make([]uuid.UUID, 0),
			}).Run()
			return &uid, err
		}
		return nil, err
	}
	return &thread.Id, nil
}

// GetThread takes the thread path and fetches it from the database
func (db *Database) GetThread(path string) (thread model.Thread, err error) {
	var result *dynamoModel.Thread

	err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbThreadTableName).Get("Path", path).One(&result)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return thread, global.ErrThreadNotFound
		}
		return thread, err
	}

	return result.ToThread(), err
}

// CreateComment takes in a body, author, and path and creates a comment for the given thread. If thread does not exist, it creates one
func (db *Database) CreateComment(body string, author string, path string, confirmed bool, replyTo *uuid.UUID) (*uuid.UUID, error) {
	thread, err := db.GetThread(path)
	if err != nil {
		if err == global.ErrThreadNotFound {
			_, err := db.CreateThread(path)
			if err != nil {
				return nil, err
			}
			return db.CreateComment(body, author, path, confirmed, replyTo)
		}
		return nil, err
	}
	if replyTo != nil {
		comment, err := db.GetComment(*replyTo)
		if err != nil {
			if err == global.ErrCommentNotFound {
				return nil, global.ErrWrongReplyTo
			}
			return nil, err
		}
		// We allow for only a single layer of nesting. (Maybe just for now? who knows.)
		// Check if the comment is a reply to this thread, and the comment you're replying to actually is a part of the thread

		if comment.ReplyTo != nil || !uuid.Equal(comment.ThreadId, thread.Id) {
			return nil, global.ErrWrongReplyTo
		}
	}
	uid := global.GetUUID()
	err = db.DB.Table(db.TablePrefix + global.DefaultDynamoDbCommentTableName).Put(model.Comment{
		Id:        uid,
		ThreadId:  thread.Id,
		Body:      body,
		Author:    author,
		Confirmed: confirmed,
		CreatedAt: time.Now().UTC(),
		ReplyTo:   replyTo,
	}).Run()
	return &uid, err
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

func (db *Database) GetUnderlyingStruct() interface{} {
	return db
}
