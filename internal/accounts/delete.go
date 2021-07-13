package accounts

import (
	"context"

	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/errors"
)

// Delete removes a user by given username and password from database
func (s Service) Delete(ctx context.Context, username string) *errors.Error {
	if err := s.Handler.Delete(ctx, "user", username); err != nil {
		return errors.CreateError(errors.DatabaseDeleteFailure, err.Error())
	}

	return nil
}
