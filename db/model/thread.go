package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// Thread represents a commenting thread
type Thread struct {
	Id        uuid.UUID `db:"Id" dynamo:"ID"`
	Path      string    `db:"Path" dynamo:"Path,hash"`
	CreatedAt time.Time `db:"CreatedAt" dynamo:"CreatedAt,range"`
}
