package authenticator

import (
	"testing"
)

func TestIncrementWithErrorAuthCounter(t *testing.T) {
	IncrementWithErrorAuthCounter("snapp", ErrIssNotFound)
}
