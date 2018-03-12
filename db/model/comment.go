package model

import "time"

// Comment represents a comment in a thread
type Comment struct {
	Id        int        `db:"Id"`
	ThreadId  int        `db:"ThreadId"`
	Body      string     `db:"Body"`
	Author    string     `db:"Author"`
	Confirmed bool       `db:"Confirmed"`
	CreatedAt time.Time  `db:"CreatedAt"`
	DeletedAt *time.Time `db:"DeletedAt"`
	ReplyTo   *int       `db:"ReplyTo"`
}
