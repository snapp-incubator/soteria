package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/accounts"
	"golang.org/x/crypto/bcrypt"
)

// Info returns a user based on given username and password
func Info(username, password string) (*User, *accounts.Error) {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return nil, accounts.CreateError(accounts.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, accounts.CreateError(accounts.WrongUsernameOrPassword, "")
	}

	return &user, nil
}
