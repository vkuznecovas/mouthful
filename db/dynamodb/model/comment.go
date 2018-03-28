package model

import (
	"time"

	"github.com/vkuznecovas/mouthful/global"

	"github.com/vkuznecovas/mouthful/db/model"

	"github.com/satori/go.uuid"
)

// Comment represents a comment in a thread
type Comment struct {
	Id        uuid.UUID `dynamo:"ID,hash"`
	ThreadId  uuid.UUID `dynamo:"ThreadId"`
	Body      string    `dynamo:"Body"`
	Author    string    `dynamo:"Author"`
	Confirmed bool      `dynamo:"Confirmed"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
	DeletedAt *int64    `dynamo:"DeletedAt,omitempty"`
	ReplyTo   *string   `dynamo:"ReplyTo,omitempty"`
}

// ToComment converts dynamoDb comment object to mouthful comment
func (c *Comment) ToComment() (model.Comment, error) {

	var deletedAt *time.Time
	if c.DeletedAt != nil {
		da := global.NanoToTime(*c.DeletedAt)
		deletedAt = &da
	}
	var replyTo *uuid.UUID
	if c.ReplyTo != nil {
		rto, err := global.ParseUUIDFromString(*c.ReplyTo)
		if err != nil {
			return model.Comment{}, err
		}
		replyTo = rto
	}
	return model.Comment{
		Id:        c.Id,
		ThreadId:  c.ThreadId,
		Body:      c.Body,
		Author:    c.Author,
		Confirmed: c.Confirmed,
		CreatedAt: c.CreatedAt,
		DeletedAt: deletedAt,
		ReplyTo:   replyTo,
	}, nil
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
