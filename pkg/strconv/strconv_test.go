package strconv_test

import (
	"testing"

	"github.com/snapp-incubator/soteria/pkg/strconv"
	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	cases := []struct {
		input  any
		output string
	}{
		{
			input:  []string{"1", "2"},
			output: "",
		},
		{
			input:  "1378",
			output: "1378",
		},
		{
			input:  1378,
			output: "1378",
		},
		{
			input:  1378.0,
			output: "1378",
		},
		{
			input:  1378.1,
			output: "1378",
		},
		{
			input:  1378.9,
			output: "1378",
		},
		{
			input:  "Hello",
			output: "Hello",
		},
	}

	for _, c := range cases {
		require.Equal(c.output, strconv.ToString(c.input))
	}
}
