package db

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

type ModelHandler interface {
	Get(pk string) user.User
}
