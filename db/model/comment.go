package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// Comment represents a comment in a thread
type Comment struct {
	Id        uuid.UUID  `db:"Id"`
	ThreadId  uuid.UUID  `db:"ThreadId"`
	Body      string     `db:"Body"`
	Author    string     `db:"Author"`
	Confirmed bool       `db:"Confirmed"`
	CreatedAt time.Time  `db:"CreatedAt"`
	DeletedAt *time.Time `db:"DeletedAt"`
	ReplyTo   *uuid.UUID `db:"ReplyTo"`
}
