package app

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/metrics"
	"sync"
)

type app struct {
	Metrics         metrics.Metrics
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

func (a *app) SetMetrics(m metrics.Metrics) {
	a.Metrics = m
}

func (a *app) SetAccountsService(service *accounts.Service) {
	a.AccountsService = service
}
