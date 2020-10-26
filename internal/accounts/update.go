package accounts

import (
	"context"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Update updates given username with given new data in `newInfo` in database
func (s Service) Update(
	ctx context.Context,
	username string,
	newPassword string,
	newType user.UserType,
	newIPs []string,
	newSecret string,
	newTokenExpiration time.Duration) *errors.Error {

	var u user.User
	if err := s.Handler.Get(context.Background(), "user", username, &u); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.CreateError(errors.PasswordHashGenerationFailure, err.Error())
		}
		u.Password = string(hash)
	}

	if newType != "" {
		u.Type = newType
	}

	if newSecret != "" {
		u.Secret = newSecret
	}

	if newIPs != nil && len(newIPs) != 0 {
		u.IPs = newIPs
	}

	if newTokenExpiration != 0 {
		u.TokenExpirationDuration = newTokenExpiration
	}

	u.MetaData.DateModified = time.Now()

	if err := s.Handler.Update(context.Background(), u); err != nil {
		return errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return nil
}
