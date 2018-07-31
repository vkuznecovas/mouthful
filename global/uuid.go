package global

import (
	"github.com/gofrs/uuid"
)

// GetUUID returns a new random UUID
func GetUUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

// ParseUUIDFromString tries to parse a string as uuid, if fails returns an error. Otherwise a pointer to uuid.UUID
func ParseUUIDFromString(uid string) (*uuid.UUID, error) {
	u2, err := uuid.FromString(uid)
	if err != nil {
		return nil, err
	}
	return &u2, nil
}
