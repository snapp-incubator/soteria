package clientid_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/clientid"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	p := clientid.NewParser(clientid.Config{
		Patterns: map[string]string{
			"android": "^android_\\d+$",
		},
	})

	require.Equal("android", p.Parse("android_124"))
	require.Equal("-", p.Parse("android_124_"))
	require.Equal("-", p.Parse("hello"))
	require.Equal("-", p.Parse("_android_124"))
}
