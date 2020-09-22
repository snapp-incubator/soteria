package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
)

// Delete removes a user by given username and password from database
func (s Service) Delete(username string) *errors.Error {
	if err := s.Handler.Delete("user", username); err != nil {
		return errors.CreateError(errors.DatabaseDeleteFailure, err.Error())
	}

	return nil
}
