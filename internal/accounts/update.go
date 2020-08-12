package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Update updates given username with given new data in `newInfo` in database
func Update(username, password, newPassword string, ips []string) *errors.Error {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.CreateError(errors.PasswordHashGenerationFailure, err.Error())
		}
		user.Password = hash
	}
	if len(ips) != 0 {
		user.IPs = ips
	}

	if err := ModelHandler.Update(user); err != nil {
		return errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return nil
}
