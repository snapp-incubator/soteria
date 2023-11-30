package authenticator

import (
	"testing"

	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/require"
)

func TestToUserAccessType(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	cases := []struct {
		name        string
		input       string
		expected    acl.AccessType
		expectedErr error
	}{
		{
			name:        "success",
			input:       "pub",
			expected:    acl.Pub,
			expectedErr: nil,
		},
		{
			name:        "failed",
			input:       "-",
			expected:    "",
			expectedErr: ErrInvalidAccessType,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			v, err := toUserAccessType(c.input)
			require.ErrorIs(c.expectedErr, err)
			require.Equal(c.expected, v)
		})
	}
}
