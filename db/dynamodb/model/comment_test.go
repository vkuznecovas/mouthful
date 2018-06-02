package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	dynamoModel "github.com/vkuznecovas/mouthful/db/dynamodb/model"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

func TestCommentConversion(t *testing.T) {
	ca := time.Now().UTC()
	da := time.Now().UTC()
	rt := global.GetUUID()
	inputMouthfulComment := model.Comment{
		Id:        global.GetUUID(),
		ThreadId:  global.GetUUID(),
		Body:      "something something",
		Author:    "Author1",
		Confirmed: true,
		CreatedAt: ca,
		DeletedAt: &da,
		ReplyTo:   &rt,
	}
	dynamoComment := dynamoModel.Comment{}
	dynamoComment.FromComment(inputMouthfulComment)
	comment, err := dynamoComment.ToComment()
	assert.Nil(t, err)
	assert.Equal(t, inputMouthfulComment.Id, comment.Id)
	assert.Equal(t, inputMouthfulComment.ThreadId, comment.ThreadId)
	assert.Equal(t, inputMouthfulComment.Body, comment.Body)
	assert.Equal(t, inputMouthfulComment.Author, comment.Author)
	assert.Equal(t, inputMouthfulComment.Confirmed, comment.Confirmed)
	assert.Equal(t, inputMouthfulComment.CreatedAt, comment.CreatedAt)
	assert.Equal(t, inputMouthfulComment.DeletedAt, comment.DeletedAt)
	assert.Equal(t, inputMouthfulComment.ReplyTo, comment.ReplyTo)
}
