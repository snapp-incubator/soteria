package db

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

type ModelHandler interface {
	Get(pk string) user.User
}

type InternalModelHanlder struct {
	users []user.User
}

func NewInternal(users []user.User) InternalModelHanlder {
	return InternalModelHanlder{
		users: users,
	}
}

func (model InternalModelHanlder) Get(pk string) user.User {
	for _, user := range model.users {
		if user.Username == pk {
			return user
		}
	}

	// nolint: exhaustivestruct
	return user.User{}
}
