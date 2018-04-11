package sqlite

import (
	"time"

	// We absolutely need the sqlite driver here, this whole package depends on it
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

var SqliteQueries = []string{
	`CREATE TABLE IF NOT EXISTS Thread(
		Id BLOB PRIMARY KEY,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null,
		Path varchar(1024) not null UNIQUE
	)`,
	`CREATE TABLE IF NOT EXISTS Comment(
		Id BLOB PRIMARY KEY,
		ThreadId INTEGER not null,
		Body text not null,
		Author varchar(255) not null,
		Confirmed bool not null default false,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null,
		ReplyTo BLOB default null,
		DeletedAt TIMESTAMP DEFAULT null,
		FOREIGN KEY(ThreadId) references Thread(Id)
	)`,
}

// InitializeDatabase runs the queries for an initial database seed
func (db *Database) InitializeDatabase() error {
	for _, v := range SqliteQueries {
		db.DB.MustExec(v)
	}
	return nil
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
	row := db.DB.QueryRowx(db.DB.Rebind("SELECT id, path, createdAt FROM Thread where path=? LIMIT 1"), path)
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
	_, err = db.DB.Exec(db.DB.Rebind("INSERT INTO comment(id, threadId, body, author, confirmed, createdAt, replyTo) VALUES(?,?,?,?,?,?,?)"), uid, thread.Id, body, author, confirmed, time.Now().UTC(), replyTo)
	return &uid, err
}

// GetCommentsByThread gets all the comments by thread path
func (db *Database) GetCommentsByThread(path string) (comments []model.Comment, err error) {
	thread, err := db.GetThread(path)
	if err != nil {
		return nil, err
	}
	err = db.DB.Select(&comments, db.DB.Rebind("select * from comment where threadId=? and confirmed=? and deletedAt is null"), thread.Id, true)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetComment gets comment by id
func (db *Database) GetComment(id uuid.UUID) (comment model.Comment, err error) {
	err = db.DB.Get(&comment, db.DB.Rebind("select * from comment where id=?"), id)
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
	res, err := db.DB.Exec(db.DB.Rebind("update comment set body=?,author=?,confirmed=? where id=?"), body, author, confirmed, id)
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
	res, err := db.DB.Exec(db.DB.Rebind("update comment set deletedAt = CURRENT_TIMESTAMP where id=? or replyTo=?"), id, id)
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
	res, err := db.DB.Exec(db.DB.Rebind("update comment set deletedAt = null where id=?"), id)
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
	err = db.DB.Select(&threads, "select * from thread")
	return threads, err
}

// GetAllComments gets all the comments found in the database
func (db *Database) GetAllComments() (comments []model.Comment, err error) {
	err = db.DB.Select(&comments, "select * from comment")
	return comments, err
}

// GetDatabaseDialect returns the current database dialect
func (db *Database) GetDatabaseDialect() string {
	return "sqlite3"
}

func (db *Database) GetUnderlyingStruct() interface{} {
	return db
}
