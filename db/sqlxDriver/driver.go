package sqlxDriver

import (
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	// We absolutely need the sqlite driver here, this whole package depends on it
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/vkuznecovas/mouthful/db/model"
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
			_, err := db.DB.Exec(db.DB.Rebind("INSERT INTO Thread(Id,Path) VALUES(?, ?)"), uid, path)
			return &uid, err
		}
		return nil, err
	}
	return &thread.Id, nil
}

// GetThread takes the thread path and fetches it from the database
func (db *Database) GetThread(path string) (thread model.Thread, err error) {
	row := db.DB.QueryRowx(db.DB.Rebind("SELECT Id, Path, CreatedAt FROM Thread where Path=? LIMIT 1"), path)
	if row != nil {
		err = row.StructScan(&thread)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return thread, global.ErrThreadNotFound
			}
			return thread, err
		}
	}
	return thread, err
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
	_, err = db.DB.Exec(db.DB.Rebind("INSERT INTO comment(Id, ThreadId, Body, Author, Confirmed, CreatedAt, ReplyTo) VALUES(?,?,?,?,?,?,?)"), uid, thread.Id, body, author, confirmed, time.Now().UTC(), replyTo)
	return &uid, err
}

// GetCommentsByThread gets all the comments by thread path
func (db *Database) GetCommentsByThread(path string) (comments []model.Comment, err error) {
	var commentSlice model.CommentSlice
	thread, err := db.GetThread(path)
	if err != nil {
		return nil, err
	}
	err = db.DB.Select(&commentSlice, db.DB.Rebind("select * from comment where ThreadId=? and Confirmed=? and DeletedAt is null"), thread.Id, true)
	if err != nil {
		return nil, err
	}
	sort.Sort(commentSlice)
	return commentSlice, nil
}

// GetComment gets comment by id
func (db *Database) GetComment(id uuid.UUID) (comment model.Comment, err error) {
	err = db.DB.Get(&comment, db.DB.Rebind("select * from comment where Id=?"), id)
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
	res, err := db.DB.Exec(db.DB.Rebind("update comment set Body=?,Author=?,Confirmed=? where Id=?"), body, author, confirmed, id)
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
	res, err := db.DB.Exec(db.DB.Rebind("update comment set DeletedAt = CURRENT_TIMESTAMP where Id=? or ReplyTo=?"), id, id)
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
	res, err := db.DB.Exec(db.DB.Rebind("update comment set DeletedAt = null where Id=?"), id)
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
	err = db.DB.Select(&threadSlice, "select * from thread")
	if err != nil {
		return threads, err
	}
	sort.Sort(threadSlice)
	return threadSlice, err
}

// GetAllComments gets all the comments found in the database
func (db *Database) GetAllComments() (comments []model.Comment, err error) {
	var commentSlice model.CommentSlice
	err = db.DB.Select(&commentSlice, "select * from comment")
	if err != nil {
		return comments, err
	}
	sort.Sort(commentSlice)
	return commentSlice, err
}

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

// WipeOutData deletes all the threads and comments in the database if the database is a test one
func (db *Database) WipeOutData() error {
	if !db.IsTest {
		return nil
	}
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		return err
	}
	_, err = tx.Exec("truncate table comment")
	if err != nil {
		return err
	}
	_, err = tx.Exec("truncate table thread")
	if err != nil {
		return err
	}
	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
