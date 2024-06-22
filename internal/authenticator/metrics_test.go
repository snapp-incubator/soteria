//nolint:testpackage
package authenticator

import (
	"testing"
)

//nolint:paralleltest
func TestIncrementWithErrorAuthCounter(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrIssNotFound)
}
