package clientid_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/clientid"
	"github.com/stretchr/testify/require"
)

func TestParser1(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	p := clientid.NewParser(clientid.Config{
		Patterns: map[string]string{
			"android": `^android_\d+$`,
		},
	})

	require.Equal("android", p.Parse("android_124"))
	require.Equal("-", p.Parse("android_124_"))
	require.Equal("-", p.Parse("hello"))
	require.Equal("-", p.Parse("_android_124"))
}

func TestParser2(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	p := clientid.NewParser(clientid.Config{
		Patterns: map[string]string{
			"android_driver": `^AD#([a-f\d]{1,32}|[A-F\d]{1,32})#\d+$`,
		},
	})

	require.Equal("android_driver", p.Parse("AD#e7adbbd4022e6ae762aae11b0593236#79971153921523"))
}
