package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// SignUp creates a user with the given information in database
func (s Service) SignUp(username, password string, userType user.UserType) *errors.Error {
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
		Password: hash,
		Type:     userType,
	}
	if err := s.Handler.Save(u); err != nil {
		return errors.CreateError(errors.DatabaseSaveFailure, err.Error())
	}

	return nil
}
