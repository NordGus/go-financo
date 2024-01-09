package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionIDLen = 64
)

var (
	ErrInvalidCredentials = errors.New("authentication: invalid credentials")
)

// authenticate receives a username and a password and authenticates the user and returns the session key. If it can't
// find the user in the system or it can't authenticated returns an error.
func (s *Service) authenticate(username string, password string) (string, error) {
	var (
		sessionBuffer = make([]byte, sessionIDLen)
		session       string
	)

	record, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return session, errors.Join(ErrInvalidCredentials, err)
	}

	err = bcrypt.CompareHashAndPassword(record.PasswordDigest(), []byte(password))
	if err != nil {
		return session, errors.Join(ErrInvalidCredentials, err)
	}

	_, err = io.ReadFull(rand.Reader, sessionBuffer)
	if err != nil {
		return session, errors.Join(ErrInvalidCredentials, err)
	}

	session = base64.URLEncoding.EncodeToString(sessionBuffer)

	err = s.sessionRepository.Create(session, record.ID())
	if err != nil {
		return session, errors.Join(ErrInvalidCredentials, err)
	}

	return session, nil
}
