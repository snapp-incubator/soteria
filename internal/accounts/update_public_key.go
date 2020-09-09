package accounts

import (
	"crypto/rsa"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"time"
)

// UpdatePublicKey updates user's public key with the given one
func (s Service) UpdatePublicKey(username string, key *rsa.PublicKey) *errors.Error {
	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	u.PublicKey = key
	u.MetaData.DateModified = time.Now()
	if err := s.Handler.Update(u); err != nil {
		return errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return nil
}
