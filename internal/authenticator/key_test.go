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
		name    string
		method  string
		keys    map[string]string
		haveErr bool
		err     error
	}{
		{
			name:   "hmac based method",
			method: "HS512",
			keys: map[string]string{
				"snpay": "YWRtaW4=",
			},
			haveErr: false,
			err:     nil,
		},
		{
			name:   "hmac based method with invalid key",
			method: "HS512",
			keys: map[string]string{
				"snpay": "YWRtaW4",
			},
			haveErr: true,
			err:     nil,
		},
		{
			name:    "invalid key type",
			method:  "Parham",
			keys:    map[string]string{},
			haveErr: true,
			err:     authenticator.ErrInvalidKeyType,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			_, err := b.GenerateKeys(c.method, c.keys)
			if c.haveErr {
				if c.err != nil {
					require.ErrorIs(t, err, c.err)
				} else {
					require.Error(t, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
