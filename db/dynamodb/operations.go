// Package dynamodb is responsible for connections and data manipulation on dynamodb
package dynamodb

import (
	"log"
	"sort"
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
		global.DefaultDynamoDbCommentTableName: dynamoModel.Comment{},
	}
	tableUnitsMap := map[string][2]int64{
		global.DefaultDynamoDbThreadTableName:  [...]int64{*db.Config.DynamoDBThreadReadUnits, *db.Config.DynamoDBThreadWriteUnits},
		global.DefaultDynamoDbCommentTableName: [...]int64{*db.Config.DynamoDBCommentReadUnits, *db.Config.DynamoDBCommentWriteUnits},
	}
	prefix := ""
	if db.Config.TablePrefix != nil {
		prefix = *db.Config.TablePrefix
	}
	db.TablePrefix = prefix
	for i := range tables {

		tables[i] = prefix + tables[i]
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
			noPrefix := strings.Replace(t, prefix, "", 1)
			readUnits := tableUnitsMap[noPrefix][0]
			writeUnits := tableUnitsMap[noPrefix][1]
			provision := db.DB.CreateTable(t, tableModelMap[noPrefix]).Provision(readUnits, writeUnits)
			if t == global.DefaultDynamoDbCommentTableName {
				provision.ProvisionIndex("ThreadId_index", *db.Config.DynamoDBIndexReadUnits, *db.Config.DynamoDBIndexWriteUnits)
			}
			err := provision.Run()
			if err != nil {
				return err
			}
		}
	}

	log.Printf("Tables created, waiting for them to be ready. Timeout - 1 minute\n")
	for i := 0; i < 60; i++ {
		dt, err := db.DB.ListTables().All()
		if err != nil {
			return err
		}

		running := 0
		for _, v := range dt {
			for _, t := range tables {
				if t == v {
					desc, err := db.DB.Table(v).Describe().Run()
					if err != nil {
						return err
					}
					if desc.Status == dynamo.ActiveStatus {
						running++
					}
				}
			}
		}
		if running == len(tables) {
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
		// Check if the comment you're replying to actually is a part of the thread
		if !uuid.Equal(comment.ThreadId, thread.Id) {
			return nil, global.ErrWrongReplyTo
		}
		// We allow for only a single layer of nesting. (Maybe just for now? who knows.)
		if comment.ReplyTo != nil && replyTo != nil {
			replyTo = comment.ReplyTo
		}
	}
	uid := global.GetUUID()
	var toReplyTo *string
	if replyTo != nil {
		trt := replyTo.String()
		toReplyTo = &trt
	}
	err = db.DB.Table(db.TablePrefix + global.DefaultDynamoDbCommentTableName).Put(dynamoModel.Comment{
		Id:        uid,
		ThreadId:  thread.Id,
		Body:      body,
		Author:    author,
		Confirmed: confirmed,
		CreatedAt: time.Now().UTC(),
		ReplyTo:   toReplyTo,
	}).Run()
	return &uid, err
}

// GetCommentsByThread gets all the comments by thread path
func (db *Database) GetCommentsByThread(path string) (comments []model.Comment, err error) {
	thread, err := db.GetThread(path)
	if err != nil {
		return comments, err
	}
	var result dynamoModel.CommentSlice

	err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Scan().Filter("'ThreadId' = ?", thread.Id).All(&result)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return comments, nil
		}
		return comments, err
	}
	sort.Sort(result)

	comments = make([]model.Comment, 0)
	for i := range result {
		comment, err := result[i].ToComment()
		if err != nil {
			return comments, err
		}
		if !comment.Confirmed {
			continue
		}
		if comment.DeletedAt != nil {
			continue
		}
		comments = append(comments, comment)
	}
	// Filter("'Count' = ? AND $ = ?", w.Count, "Message", w.Msg)
	return comments, nil
}

// GetComment gets comment by id
func (db *Database) GetComment(id uuid.UUID) (comment model.Comment, err error) {
	var result *dynamoModel.Comment

	err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Get("ID", id).One(&result)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return comment, global.ErrCommentNotFound
		}
		return comment, err
	}
	res, err := result.ToComment()
	return res, err
}

// UpdateComment updatesComment comment by id
func (db *Database) UpdateComment(id uuid.UUID, body, author string, confirmed bool) error {
	_, err := db.GetComment(id)
	if err != nil {
		return err
	}
	statement := db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Update("ID", id)
	statement.Set("Body", body)
	statement.Set("Author", author)
	statement.Set("Confirmed", confirmed)
	err = statement.Run()
	return err
}

// DeleteComment soft-deletes the comment by id and all the replies to it
func (db *Database) DeleteComment(id uuid.UUID) error {
	comment, err := db.GetComment(id)
	if err != nil {
		return err
	}

	deletedAt := time.Now().UnixNano()
	err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Update("ID", id).Set("DeletedAt", deletedAt).Run()
	if err != nil {
		return err
	}
	var result dynamoModel.CommentSlice
	err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Scan().Filter("'ThreadId' = ?", comment.ThreadId).All(&result)
	if err != nil {
		return err
	}
	for i := range result {
		if result[i].ReplyTo != nil {
			cid, err := global.ParseUUIDFromString(*result[i].ReplyTo)
			if err != nil {
				return err
			}
			if uuid.Equal(*cid, comment.Id) {
				err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Update("ID", result[i].Id).Set("DeletedAt", deletedAt).Run()
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

// RestoreDeletedComment restores the soft-deleted comment
func (db *Database) RestoreDeletedComment(id uuid.UUID) error {
	_, err := db.GetComment(id)
	if err != nil {
		return err
	}
	err = db.DB.Table(db.TablePrefix+global.DefaultDynamoDbCommentTableName).Update("ID", id).Remove("DeletedAt").Run()
	return err
}

// GetAllThreads gets all the threads found in the database
func (db *Database) GetAllThreads() (threads []model.Thread, err error) {
	var result dynamoModel.ThreadSlice
	db.DB.Table(db.TablePrefix + global.DefaultDynamoDbThreadTableName).Scan().All(&result)
	sort.Sort(result)
	threads = make([]model.Thread, len(result))
	for i := range result {
		threads[i] = result[i].ToThread()
	}
	return threads, err
}

// GetAllComments gets all the comments found in the database
func (db *Database) GetAllComments() (comments []model.Comment, err error) {
	var result dynamoModel.CommentSlice
	db.DB.Table(db.TablePrefix + global.DefaultDynamoDbCommentTableName).Scan().All(&result)
	sort.Sort(result)
	comments = make([]model.Comment, len(result))
	for i := range result {
		comment, err := result[i].ToComment()
		if err != nil {
			return comments, err
		}
		comments[i] = comment
	}
	return comments, err
}

// GetDatabaseDialect returns the current database dialect
func (db *Database) GetDatabaseDialect() string {
	return "dynamodb"
}

// GetUnderlyingStruct returns the underlying database struct for the driver
func (db *Database) GetUnderlyingStruct() interface{} {
	return db
}
