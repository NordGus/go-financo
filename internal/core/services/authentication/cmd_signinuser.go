package authentication

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"financo/internal/core/domain/session"
	"financo/internal/core/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type SignInUserReq struct {
	Username string
	Password string
}

type SignInUserResp struct {
	SessionID string
}

func (s *service) SignIn(ctx context.Context, req SignInUserReq) (SignInUserResp, error) {
	sessionIDBuf := make([]byte, session.IdByteSize)

	username, err := user.NewUsername(req.Username)
	if err != nil {
		return SignInUserResp{}, err
	}

	usr, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return SignInUserResp{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordDigest.String()), []byte(req.Password))
	if err != nil {
		return SignInUserResp{}, err
	}

	_, err = io.ReadFull(rand.Reader, sessionIDBuf)
	if err != nil {
		return SignInUserResp{}, err
	}

	sessionID, err := session.NewID([session.IdByteSize]byte(sessionIDBuf))
	if err != nil {
		return SignInUserResp{}, err
	}

	err = s.sessionRepo.Create(ctx, session.New(sessionID, usr.ID, usr.SessionVersion, time.Now().UTC()))
	if err != nil {
		return SignInUserResp{}, err
	}

	return SignInUserResp{SessionID: sessionID.String()}, nil
}
