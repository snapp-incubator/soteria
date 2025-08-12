package metric

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

func register[T prometheus.Collector](metric T) T {
	if err := prometheus.Register(metric); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &are); ok {
			metric, ok = are.ExistingCollector.(T)
			if !ok {
				panic("different metric type registered with the same name")
			}
		} else {
			panic(err)
		}
	}

	return metric
}
