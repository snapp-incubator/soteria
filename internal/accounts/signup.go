package accounts

import (
	"context"
	"time"

	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

// SignUp creates a user with the given information in database
func (s Service) SignUp(ctx context.Context, username, password string, userType user.Type) *errors.Error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.CreateError(errors.PasswordHashGenerationFailure, err.Error())
	}

	u := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username: username,
		Password: string(hash),
		Type:     userType,
	}
	if err := s.Handler.Save(context.Background(), u); err != nil {
		return errors.CreateError(errors.DatabaseSaveFailure, err.Error())
	}

	return nil
}
