package app

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	Authenticator *authenticator.Authenticator
	Tracer        trace.Tracer
}
