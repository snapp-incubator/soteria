package app

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
	"sync"
)

type app struct {
	Authenticator *accounts.Authenticator
}

var singleton *app
var once sync.Once

func GetInstance() *app {
	once.Do(func() {
		singleton = &app{}
	})
	return singleton
}

func (a *app) SetAuthenticator(authenticator *accounts.Authenticator) {
	a.Authenticator = authenticator
}
