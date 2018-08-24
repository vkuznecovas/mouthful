// Package model contains the models consumed by the databases.
package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Comment represents a comment in a thread
type Comment struct {
	Id        uuid.UUID  `db:"Id" json:"Id"`
	ThreadId  uuid.UUID  `db:"ThreadId" json:"ThreadId"`
	Body      string     `db:"Body" json:"Body"`
	Author    string     `db:"Author" json:"Author"`
	Confirmed bool       `db:"Confirmed" json:"Confirmed"`
	CreatedAt time.Time  `db:"CreatedAt" json:"CreatedAt"`
	DeletedAt *time.Time `db:"DeletedAt" json:"DeletedAt,omitempty"`
	ReplyTo   *uuid.UUID `db:"ReplyTo" json:"ReplyTo,omitempty"`
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
