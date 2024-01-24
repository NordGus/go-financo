package session

import (
	"errors"
	"time"

	"github.com/NordGus/fnncr/internal/core/domain/user"
	"github.com/google/uuid"
)

var (
	ErrExpired = errors.New("session: session expired")
)

type Session struct {
	ID        ID
	UserID    uuid.UUID
	Version   int32
	CreatedAt time.Time
}

func New(id ID, userId uuid.UUID, version int32, ca time.Time) Session {
	return Session{
		ID:        id,
		UserID:    userId,
		Version:   version,
		CreatedAt: ca,
	}
}

func (s Session) Expired(usr user.User) bool {
	// TODO: implement a session lifespan
	return s.Version != usr.SessionVersion
}
