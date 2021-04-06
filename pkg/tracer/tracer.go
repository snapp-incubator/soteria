package tracer

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	jaegerConf "github.com/uber/jaeger-client-go/config"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
)

// New receives a `configs.TracerConfig` and returns a `opentracing.Tracer` and a `io.Closer` and a `error` if there was one
func New(cfg *config.TracerConfig) (opentracing.Tracer, io.Closer, error) {
	trc, cl, err := jaegerConf.Configuration{
		ServiceName: cfg.ServiceName,
		Disabled:    !cfg.Enabled,
		Sampler: &jaegerConf.SamplerConfig{
			Type:  cfg.SamplerType,
			Param: cfg.SamplerParam,
		},
		Reporter: &jaegerConf.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
	}.NewTracer()

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new tracer: %w", err)
	}

	return trc, cl, nil
}
