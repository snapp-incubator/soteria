package app

import (
	"io"
	"sync"

	"github.com/opentracing/opentracing-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/emq"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/metrics"
)

type app struct {
	Metrics         metrics.Metrics
	Authenticator   *authenticator.Authenticator
	AccountsService *accounts.Service
	EMQStore        emq.Store
	Tracer          opentracing.Tracer
	TracerCloser    io.Closer
}

var (
	singleton *app
	once      sync.Once
)

func GetInstance() *app {
	once.Do(func() {
		singleton = &app{}
	})

	return singleton
}

func (a *app) SetEMQStore(store emq.Store) {
	a.EMQStore = store
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

func (a *app) SetTracer(tracer opentracing.Tracer, closer io.Closer) {
	a.Tracer = tracer
	a.TracerCloser = closer
}
