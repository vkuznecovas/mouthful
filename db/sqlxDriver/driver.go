// Package sqlxDriver allows for data manipulation through the github.com/jmoiron/sqlx package.
package sqlxDriver

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	// We absolutely need the sqlite driver here, this whole package depends on it
	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vkuznecovas/mouthful/db/model"
	tool "github.com/vkuznecovas/mouthful/db/tool"

	"github.com/vkuznecovas/mouthful/global"
)

// Database is a database instance for sqlx
type Database struct {
	DB      *sqlx.DB
	Queries []string
	Dialect string
	IsTest  bool
}

// CreateThread takes the thread path and creates it in the database
func (db *Database) CreateThread(path string) (*uuid.UUID, error) {
	thread, err := db.GetThread(path)
	if err != nil {
		if err == global.ErrThreadNotFound {
			uid := global.GetUUID()
			res, err := db.DB.Exec(db.DB.Rebind("INSERT INTO Thread(Id,Path) VALUES(?, ?)"), uid, path)
			if err != nil {
				return nil, err
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, err
			}
			if affected != 1 {
				return nil, global.ErrInternalServerError
			}
			return &uid, err
		}
		return nil, err
	}
	return &thread.Id, nil
}

// GetThread takes the thread path and fetches it from the database
func (db *Database) GetThread(path string) (thread model.Thread, err error) {
	err = db.DB.QueryRowx(db.DB.Rebind("SELECT Id, Path, CreatedAt FROM Thread where Path=? LIMIT 1"), path).StructScan(&thread)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return thread, global.ErrThreadNotFound
		}
		return thread, err
	}
	return thread, err
}

// CreateComment takes in a body, author, and path and creates a comment for the given thread. If thread does not exist, it creates one
func (db *Database) CreateComment(body string, author string, path string, confirmed bool, replyTo *uuid.UUID) (*uuid.UUID, error) {
	thread, err := db.GetThread(path)
	if err != nil {
		if err == global.ErrThreadNotFound {
			threadId, err := db.CreateThread(path)
			if err != nil {
				return nil, err
			}
			uid := global.GetUUID()
			if replyTo != nil {
				return nil, global.ErrWrongReplyTo
			}
			res, err := db.DB.Exec(db.DB.Rebind("INSERT INTO Comment(Id, ThreadId, Body, Author, Confirmed, CreatedAt, ReplyTo) VALUES(?,?,?,?,?,?,?)"), uid, threadId, body, author, confirmed, time.Now().UTC(), nil)
			if err != nil {
				return nil, err
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, err
			}
			if affected != 1 {
				return nil, global.ErrInternalServerError
			}
			return &uid, err
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
		if !bytes.Equal(comment.ThreadId.Bytes(), thread.Id.Bytes()) {
			return nil, global.ErrWrongReplyTo
		}
		// We allow for only a single layer of nesting. (Maybe just for now? who knows.)
		if comment.ReplyTo != nil && replyTo != nil {
			replyTo = comment.ReplyTo
		}
	}
	uid := global.GetUUID()
	res, err := db.DB.Exec(db.DB.Rebind("INSERT INTO Comment(Id, ThreadId, Body, Author, Confirmed, CreatedAt, ReplyTo) VALUES(?,?,?,?,?,?,?)"), uid, thread.Id, body, author, confirmed, time.Now().UTC(), replyTo)
	if err != nil {
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &uid, err
}

// GetCommentsByThread gets all the comments by thread path
func (db *Database) GetCommentsByThread(path string) (comments []model.Comment, err error) {
	var commentSlice model.CommentSlice
	thread, err := db.GetThread(path)
	if err != nil {
		return nil, err
	}
	err = db.DB.Select(&commentSlice, db.DB.Rebind("select * from Comment where ThreadId=? and Confirmed=? and DeletedAt is null"), thread.Id, true)
	if err != nil {
		return nil, err
	}
	sort.Sort(commentSlice)
	return commentSlice, nil
}

// GetComment gets comment by id
func (db *Database) GetComment(id uuid.UUID) (comment model.Comment, err error) {
	err = db.DB.Get(&comment, db.DB.Rebind("select * from Comment where Id=?"), id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return comment, global.ErrCommentNotFound
		}
		return comment, err
	}
	return comment, nil
}

// UpdateComment updatesComment comment by id
func (db *Database) UpdateComment(id uuid.UUID, body, author string, confirmed bool) error {
	res, err := db.DB.Exec(db.DB.Rebind("update Comment set Body=?,Author=?,Confirmed=? where Id=?"), body, author, confirmed, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return global.ErrCommentNotFound
	}
	return nil
}

// DeleteComment soft-deletes the comment by id and all the replies to it
func (db *Database) DeleteComment(id uuid.UUID) error {
	res, err := db.DB.Exec(db.DB.Rebind("update Comment set DeletedAt = CURRENT_TIMESTAMP where Id=? or ReplyTo=?"), id, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return global.ErrCommentNotFound
	}
	return nil
}

// RestoreDeletedComment restores the soft-deleted comment
func (db *Database) RestoreDeletedComment(id uuid.UUID) error {
	res, err := db.DB.Exec(db.DB.Rebind("update Comment set DeletedAt = null where Id=?"), id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return global.ErrCommentNotFound
	}
	return nil
}

// GetAllThreads gets all the threads found in the database
func (db *Database) GetAllThreads() (threads []model.Thread, err error) {
	var threadSlice model.ThreadSlice
	err = db.DB.Select(&threadSlice, "select * from Thread")
	if err != nil {
		return threads, err
	}
	sort.Sort(threadSlice)
	return threadSlice, err
}

// GetAllComments gets all the comments found in the database
func (db *Database) GetAllComments() (comments []model.Comment, err error) {
	var commentSlice model.CommentSlice
	err = db.DB.Select(&commentSlice, "select * from Comment")
	if err != nil {
		return comments, err
	}
	sort.Sort(commentSlice)
	return commentSlice, err
}

// GetUnderlyingStruct returns the underlying database struct for the driver
func (db *Database) GetUnderlyingStruct() interface{} {
	return db
}

// InitializeDatabase runs the queries for an initial database seed
func (db *Database) InitializeDatabase() error {
	for _, v := range db.Queries {
		db.DB.MustExec(v)
	}
	return nil
}

// GetDatabaseDialect returns the current database dialect
func (db *Database) GetDatabaseDialect() string {
	return db.Dialect
}

// CleanUpStaleData removes the stale data from the database
func (db *Database) CleanUpStaleData(target global.CleanupType, timeout int64) error {
	timeoutDuration := time.Duration(int64(time.Second) * timeout)
	deleteFrom := time.Now().Add(-timeoutDuration).UTC()
	if target == global.Deleted {
		return db.CleanupDeleted(deleteFrom)
	} else if target == global.Unconfirmed {
		return db.CleanupUnconfirmed(deleteFrom)
	}
	return fmt.Errorf("Unknown cleanup type %v", target)
}

// HardDeleteComment permanently deletes the comment from a database.
func (db *Database) HardDeleteComment(commentId uuid.UUID) error {
	res, err := db.DB.Exec(db.DB.Rebind("delete from Comment where Id=? or ReplyTo=?"), commentId, commentId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return global.ErrCommentNotFound
	}
	return nil
}

// CleanupUnconfirmed removes the unconfirmed comments that are older than the given time
func (db *Database) CleanupUnconfirmed(olderThan time.Time) error {
	query := db.DB.Rebind("select * from Comment where Confirmed=? and DeletedAt is null")
	var commentSlice model.CommentSlice
	err := db.DB.Select(&commentSlice, query, false)
	if err != nil {
		return err
	}
	for _, v := range commentSlice {
		if v.CreatedAt.Before(olderThan) {
			err = db.HardDeleteComment(v.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanupDeleted removes the deleted comments that are older than the given time
func (db *Database) CleanupDeleted(olderThan time.Time) error {
	query := "select * from Comment where DeletedAt is not null"
	var commentSlice model.CommentSlice
	err := db.DB.Select(&commentSlice, query)
	if err != nil {
		return err
	}
	for _, v := range commentSlice {
		if v.DeletedAt.Before(olderThan) {
			err = db.HardDeleteComment(v.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WipeOutData deletes all the threads and comments in the database if the database is a test one
func (db *Database) WipeOutData() error {
	if !db.IsTest {
		return nil
	}
	if db.Dialect == "postgres" {
		_, err := db.DB.Exec("truncate table Thread CASCADE")
		if err != nil {
			return err
		}
		return nil
	}
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	if db.Dialect == "mysql" {
		_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0")
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec("truncate table Comment")
	if err != nil {
		return err
	}
	_, err = tx.Exec("truncate table Thread")
	if err != nil {
		return err
	}
	if db.Dialect == "mysql" {
		_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// ImportData performs the data import for the given driver
func (db *Database) ImportData(pathToDump string) error {
	importThread := func(t model.Thread) error {
		_, err := db.DB.Exec(db.DB.Rebind("INSERT INTO Thread(Id,Path,CreatedAt) VALUES(?, ?, ?)"), t.Id, t.Path, t.CreatedAt)
		if err != nil {
			return err
		}
		return nil
	}
	importComment := func(c model.Comment) error {
		_, err := db.DB.Exec(db.DB.Rebind("INSERT INTO Comment(Id, ThreadId, Body, Author, Confirmed, CreatedAt, ReplyTo, DeletedAt) VALUES(?,?,?,?,?,?,?,?)"), c.Id, c.ThreadId, c.Body, c.Author, c.Confirmed, c.CreatedAt, c.ReplyTo, c.DeletedAt)
		if err != nil {
			return err
		}
		return nil
	}
	err := tool.ImportData(pathToDump, importThread, importComment)
	return err
}
