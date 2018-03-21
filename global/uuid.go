package global

import (
	"github.com/satori/go.uuid"
)

func GetUUID() uuid.UUID {
	return uuid.NewV4()
}

func ParseUUIDFromString(uid string) (*uuid.UUID, error) {
	u2, err := uuid.FromString(uid)
	if err != nil {
		return nil, err
	}
	return &u2, nil
}
