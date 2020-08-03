package accounts

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// SignUp is for creating new users
func SignUp(username, password, userType string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	if err != nil {
		return err
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
	return ModelHandler.Save(user)
}
