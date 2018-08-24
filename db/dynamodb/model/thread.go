// Package model deals with specific models for dynamodb(they might be a bit different to db/model). It allows for sorting of the given models and transformation to db/model equivalents.
package model

import (
	"time"

	"github.com/vkuznecovas/mouthful/db/model"

	"github.com/gofrs/uuid"
)

// Thread represents a commenting thread for dynamodb.
// This is needed since we store threads a bit differently than in a relational database
type Thread struct {
	Id        uuid.UUID `dynamo:"ID"`
	Path      string    `dynamo:"Path,hash"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// ToThread converts dynamodb thread to mouthful thread
func (t *Thread) ToThread() model.Thread {
	return model.Thread{
		Id:        t.Id,
		Path:      t.Path,
		CreatedAt: t.CreatedAt,
	}
}

// ThreadSlice represents a collection of threads
type ThreadSlice []Thread

func (ts ThreadSlice) Len() int {
	return len(ts)
}

func (ts ThreadSlice) Less(i, j int) bool {
	return ts[i].CreatedAt.Before(ts[j].CreatedAt)
}

func (ts ThreadSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}
