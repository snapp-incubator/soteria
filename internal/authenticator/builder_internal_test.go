package authenticator

import (
	"testing"

	"github.com/snapp-incubator/soteria/pkg/acl"
	"github.com/stretchr/testify/require"
)

// nolint: funlen
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
			name:        "success",
			input:       "publish",
			expected:    acl.Pub,
			expectedErr: nil,
		},
		{
			name:        "success",
			input:       "sub",
			expected:    acl.Sub,
			expectedErr: nil,
		},
		{
			name:        "success",
			input:       "subscribe",
			expected:    acl.Sub,
			expectedErr: nil,
		},
		{
			name:        "success",
			input:       "pubsub",
			expected:    acl.PubSub,
			expectedErr: nil,
		},
		{
			name:        "success",
			input:       "subpub",
			expected:    acl.PubSub,
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
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			v, err := toUserAccessType(c.input)
			require.ErrorIs(c.expectedErr, err)
			require.Equal(c.expected, v)
		})
	}
}
