package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

// Delete removes a user by given username and password from database
func (s Service) Delete(username, password string) *errors.Error {
	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	if err := s.Handler.Delete("user", username); err != nil {
		return errors.CreateError(errors.DatabaseDeleteFailure, err.Error())
	}

	return nil
}
