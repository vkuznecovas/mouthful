package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// Thread represents a commenting thread
type Thread struct {
	Id        uuid.UUID `db:"Id"`
	Path      string    `db:"Path"`
	CreatedAt time.Time `db:"CreatedAt"`
}
