package authenticator_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeys(t *testing.T) {
	t.Parallel()

	b := new(authenticator.Builder)

	cases := []struct {
		name   string
		method string
		keys   map[string]string
		err    error
	}{
		{
			name:   "hmac based method",
			method: "HS512",
			keys: map[string]string{
				"snpay": "YWRtaW4=",
			},
			err: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			_, err := b.GenerateKeys(c.method, c.keys)
			if c.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, c.err)
			}
		})
	}
}
