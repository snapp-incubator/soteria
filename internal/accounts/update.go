package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

// Update updates given username with given new data in `newInfo` in database
func (s Service) Update(username, password, newPassword, secret string, ips []string) *errors.Error {
	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		return errors.CreateError(errors.WrongUsernameOrPassword, "")
	}

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.CreateError(errors.PasswordHashGenerationFailure, err.Error())
		}
		u.Password = hash
	}

	if secret != "" {
		u.Secret = secret
	}

	if len(ips) != 0 {
		u.IPs = ips
	}

	if err := s.Handler.Update(u); err != nil {
		return errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return nil
}
