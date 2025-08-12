package profiler

import (
	"log"
	"os"

	"github.com/grafana/pyroscope-go"
)

func Start(cfg Config) {
	if cfg.Enabled {
		// nolint: exhaustruct
		if _, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: "snapp.soteria",
			ServerAddress:   cfg.URL,
			Tags:            map[string]string{"hostname": os.Getenv("HOSTNAME")},

			ProfileTypes: []pyroscope.ProfileType{
				// these profile types are enabled by default:
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,

				// these profile types are optional:
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
				pyroscope.ProfileBlockCount,
				pyroscope.ProfileBlockDuration,
			},
		}); err != nil {
			log.Printf("failed to start the profiler")
		}
	}
}
