package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/accounts"
	"golang.org/x/crypto/bcrypt"
	"net"
)

// Update updates given username with given new data in `newInfo` in database
func Update(username, password, newPassword string, ips []net.IP) *accounts.Error {
	var user User
	if err := ModelHandler.Get("user", username, &user); err != nil {
		return accounts.CreateError(accounts.DatabaseGetFailure, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return accounts.CreateError(accounts.WrongUsernameOrPassword, "")
	}

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return accounts.CreateError(accounts.PasswordHashGenerationFailure, err.Error())
		}
		user.Password = hash
	}
	if len(ips) != 0 {
		user.IPs = ips
	}

	if err := ModelHandler.Update(user); err != nil {
		return accounts.CreateError(accounts.DatabaseUpdateFailure, err.Error())
	}

	return nil
}
