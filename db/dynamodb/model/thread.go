package model

import (
	"time"

	"github.com/vkuznecovas/mouthful/db/model"

	"github.com/satori/go.uuid"
)

// Thread represents a commenting thread for dynamodb.
// This is needed since we store threads a bit differently than in a relational database
type Thread struct {
	Id        uuid.UUID   `db:"Id" dynamo:"ID"`
	Path      string      `db:"Path" dynamo:"Path,hash"`
	CreatedAt time.Time   `db:"CreatedAt" dynamo:"CreatedAt"`
	Comments  []uuid.UUID `dynamo:"Comments"`
}

// ToThread converts dynamodb thread to mouthful thread
func (t *Thread) ToThread() model.Thread {
	return model.Thread{
		Id:        t.Id,
		Path:      t.Path,
		CreatedAt: t.CreatedAt,
	}
}
