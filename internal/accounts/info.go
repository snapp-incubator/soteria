package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

// Info returns a user based on given username and password
func (s Service) Info(username, password string) (*user.User, *errors.Error) {
	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return nil, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	return &u, nil
}
