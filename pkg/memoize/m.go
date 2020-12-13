package memoize

import (
	"golang.org/x/crypto/bcrypt"
)

func MemoizedCompareHashAndPassword() func([]byte, []byte) error {
	type args struct {
		hashedPassword string
		password       string
	}
	cache := make(map[args]error)
	return func(hashedPassword []byte, password []byte) error {
		key := args{
			hashedPassword: string(hashedPassword),
			password:       string(password),
		}
		if _, ok := cache[key]; !ok {
			cache[key] = bcrypt.CompareHashAndPassword(hashedPassword, password)
		}
		return cache[key]
	}
}
