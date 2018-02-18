package sqlite

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

var sqliteQueries = []string{
	`CREATE TABLE IF NOT EXISTS Thread(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Path varchar(1024) not null UNIQUE
	)`,
	`CREATE TABLE IF NOT EXISTS Comment(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		ThreadId INTEGER not null,
		Body text not null,
		Author varchar(255) not null,
		Confirmed bool not null default false,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null,
		ReplyTo INTEGER default null,
		FOREIGN KEY(ThreadId) references Thread(Id)
	)`,
}

func (db *Database) InitializeDatabase() error {
	for _, v := range sqliteQueries {
		db.DB.MustExec(v)
	}
	return nil
}

// CreateThread takes the thread path and creates it in the database
func (db *Database) CreateThread(path string) error {
	_, err := db.DB.Exec(db.DB.Rebind("INSERT INTO thread(Path) VALUES(?)"), path)
	return err
}

// GetThread takes the thread path and fetches it from the database
func (db *Database) GetThread(path string) (thread model.Thread, err error) {
	row := db.DB.QueryRowx(db.DB.Rebind("SELECT id, path FROM thread where path=? LIMIT 1"), path)
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
func (db *Database) CreateComment(body string, author string, path string, confirmed bool, replyTo *int) error {
	thread, err := db.GetThread(path)
	if err != nil {
		if err == global.ErrThreadNotFound {
			err = db.CreateThread(path)
			if err != nil {
				return err
			}
			return db.CreateComment(body, author, path, confirmed, replyTo)
		}
		return err
	}
	if replyTo != nil {
		comment, err := db.GetComment(*replyTo)
		if err != nil {
			return err
		}
		// We allow for only a single layer of nesting. (Maybe just for now? who knows.)
		// Check if the comment is a reply to this thread, and the comment you're replying to actually is a part of the thread
		if comment.ReplyTo != nil || comment.ThreadId != thread.Id {
			return global.ErrWrongReplyTo
		}
	}
	_, err = db.DB.Exec(db.DB.Rebind("INSERT INTO comment(threadId, body, author, confirmed, createdAt, replyTo) VALUES(?,?,?,?,?,?)"), thread.Id, body, author, confirmed, time.Now().UTC(), replyTo)
	return err
}

// GetCommentsByThread gets all the comments by thread path
func (db *Database) GetCommentsByThread(path string) (comments []model.Comment, err error) {
	thread, err := db.GetThread(path)
	if err != nil {
		return nil, err
	}
	err = db.DB.Select(&comments, db.DB.Rebind("select * from comment where threadId=? and confirmed=?"), thread.Id, true)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetComment gets comment by id
func (db *Database) GetComment(id int) (comment model.Comment, err error) {
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
func (db *Database) UpdateComment(id int, body, author string, confirmed bool) error {
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

// DeleteComment deletes the comment by id
func (db *Database) DeleteComment(id int) error {
	res, err := db.DB.Exec(db.DB.Rebind("delete from comment where id=?"), id)
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
