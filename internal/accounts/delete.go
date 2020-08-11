package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/accounts"
	"golang.org/x/crypto/bcrypt"
)

// Delete removes a user by given username and password from database
func Delete(username, password string) *accounts.Error {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return accounts.CreateError(accounts.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return accounts.CreateError(accounts.WrongUsernameOrPassword, "")
	}

	if err := ModelHandler.Delete("user", username); err != nil {
		return accounts.CreateError(accounts.DatabaseDeleteFailure, err.Error())
	}

	return nil
}
