package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// SignUp creates a user with the given information in database
func SignUp(username, password, userType string) *errors.Error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.CreateError(errors.PasswordHashGenerationFailure, err.Error())
	}

	user := User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username: username,
		Password: hash,
		Type:     userType,
	}
	if err := ModelHandler.Save(user); err != nil {
		return errors.CreateError(errors.DatabaseSaveFailure, err.Error())
	}

	return nil
}
