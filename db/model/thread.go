package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Thread represents a commenting thread
type Thread struct {
	Id        uuid.UUID `db:"Id" dynamo:"ID" json:"Id"`
	Path      string    `db:"Path" dynamo:"Path,hash" json:"Path"`
	CreatedAt time.Time `db:"CreatedAt" dynamo:"CreatedAt,range" json:"CreatedAt,omitempty"`
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
