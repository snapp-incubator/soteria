package app

import (
	"io"
	"sync"

	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	Authenticator *authenticator.Authenticator
	Tracer        trace.Tracer
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

func (a *App) SetTracer(tracer trace.Tracer) {
	a.Tracer = tracer
}
