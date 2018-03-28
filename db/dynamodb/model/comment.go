package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// Comment represents a comment in a thread
type Comment struct {
	Id        uuid.UUID  `dynamo:"ID,hash"`
	ThreadId  uuid.UUID  `dynamo:"ThreadId"`
	Body      string     `dynamo:"Body"`
	Author    string     `dynamo:"Author"`
	Confirmed bool       `dynamo:"Confirmed"`
	CreatedAt time.Time  `dynamo:"CreatedAt"`
	DeletedAt float64    `dynamo:"DeletedAt,omitempty"`
	ReplyTo   *uuid.UUID `dynamo:"ReplyTo,omitempty"`
}
