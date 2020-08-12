package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Info returns a user based on given username and password
func Info(username, password string) (*User, *errors.Error) {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return nil, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	return &user, nil
}
