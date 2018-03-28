package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// Comment represents a comment in a thread
type Comment struct {
	Id        uuid.UUID  `db:"Id" dynamo:"ID,hash"`
	ThreadId  uuid.UUID  `db:"ThreadId"  dynamo:"ThreadId"`
	Body      string     `db:"Body"  dynamo:"Body"`
	Author    string     `db:"Author"  dynamo:"Author"`
	Confirmed bool       `db:"Confirmed"  dynamo:"Confirmed"`
	CreatedAt time.Time  `db:"CreatedAt" dynamo:"CreatedAt"`
	DeletedAt *time.Time `db:"DeletedAt" dynamo:"DeletedAt,omitempty"`
	ReplyTo   *uuid.UUID `db:"ReplyTo"  dynamo:"ReplyTo,omitempty"`
}
