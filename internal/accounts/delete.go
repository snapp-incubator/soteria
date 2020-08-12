package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Delete removes a user by given username and password from database
func Delete(username, password string) *errors.Error {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	if err := ModelHandler.Delete("user", username); err != nil {
		return errors.CreateError(errors.DatabaseDeleteFailure, err.Error())
	}

	return nil
}
