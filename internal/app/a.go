package app

import (
	"io"
	"sync"

	"github.com/opentracing/opentracing-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/metrics"
)

type App struct {
	Metrics       metrics.Metrics
	Authenticator *authenticator.Authenticator
	Tracer        opentracing.Tracer
	TracerCloser  io.Closer
}

// nolint: gochecknoglobals
var (
	singleton *App
	once      sync.Once
)

func GetInstance() *App {
	once.Do(func() {
		singleton = new(App)
	})

	return singleton
}

func (a *App) SetAuthenticator(authenticator *authenticator.Authenticator) {
	a.Authenticator = authenticator
}

func (a *App) SetMetrics(m metrics.Metrics) {
	a.Metrics = m
}

func (a *App) SetTracer(tracer opentracing.Tracer, closer io.Closer) {
	a.Tracer = tracer
	a.TracerCloser = closer
}
