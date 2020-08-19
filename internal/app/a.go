package app

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"sync"
)

type app struct {
	Authenticator   *authenticator.Authenticator
	AccountsService *accounts.Service
}

var singleton *app
var once sync.Once

func GetInstance() *app {
	once.Do(func() {
		singleton = &app{}
	})
	return singleton
}

func (a *app) SetAuthenticator(authenticator *authenticator.Authenticator) {
	a.Authenticator = authenticator
}

func (a *app) SetAccountsService(service *accounts.Service) {
	a.AccountsService = service
}
