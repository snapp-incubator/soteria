package accounts

import (
	"context"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

// Info returns a user based on given username and password
func (s Service) Info(ctx context.Context, username, password string) (*user.User, *errors.Error) {
	var u user.User
	if err := s.Handler.Get(ctx, "user", username, &u); err != nil {
		return nil, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	return &u, nil
}
