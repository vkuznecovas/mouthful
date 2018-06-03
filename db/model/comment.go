// Package model contains the models consumed by the databases.
package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// Comment represents a comment in a thread
type Comment struct {
	Id        uuid.UUID  `db:"Id" json:"id"`
	ThreadId  uuid.UUID  `db:"ThreadId" json:"threadid"`
	Body      string     `db:"Body" json:"body"`
	Author    string     `db:"Author" json:"author"`
	Confirmed bool       `db:"Confirmed" json:"confirmed"`
	CreatedAt time.Time  `db:"CreatedAt" json:"createdAt"`
	DeletedAt *time.Time `db:"DeletedAt" json:"deletedAt,omitempty"`
	ReplyTo   *uuid.UUID `db:"ReplyTo" json:"replyTo,omitempty"`
}

// CommentSlice represents a collection of comments
type CommentSlice []Comment

func (cs CommentSlice) Len() int {
	return len(cs)
}

func (cs CommentSlice) Less(i, j int) bool {
	return cs[i].CreatedAt.Before(cs[j].CreatedAt)
}

func (cs CommentSlice) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
