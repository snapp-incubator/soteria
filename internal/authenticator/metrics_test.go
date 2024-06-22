package authenticator

import (
	"testing"
)

func TestIncrementWithErrorAuthCounter(_ *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrIssNotFound)
}
