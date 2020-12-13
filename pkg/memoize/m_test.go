package memoize

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestMemoizedCompareHashAndPassword(t *testing.T) {
	fn := MemoizedCompareHashAndPassword()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	actual := fn(hashedPassword, []byte("test"))
	expected := bcrypt.CompareHashAndPassword(hashedPassword, []byte("test"))
	assert.Equal(t, actual, expected)

	s1 :=  time.Now()
	for i := 0; i < 10; i ++ {
		_ = bcrypt.CompareHashAndPassword(hashedPassword, []byte("test"))
	}
	d1 := time.Since(s1)

	s2 :=  time.Now()
	for i := 0; i < 10; i ++ {
		_ = fn(hashedPassword, []byte("test"))
	}
	d2 := time.Since(s2)
	assert.Less(t, d2.Nanoseconds(), d1.Nanoseconds())
}
