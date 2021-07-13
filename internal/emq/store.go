package emq

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Store struct {
	Client *redis.Client
}

var (
	ErrInvalidHResult = errors.New("invalid hresult structure")
	ErrUserNotExists  = errors.New("user doesn't exist")
)

const (
	PasswordHKey    = "password"
	IsSuperuserHKey = "is_superuser"
)

func Key(username string) string {
	return fmt.Sprintf("mqtt_user:%s", username)
}

func (s *Store) Save(ctx context.Context, u User) error {
	if err := s.Client.HSet(ctx, Key(u.Username), PasswordHKey, u.Password,
		IsSuperuserHKey, strconv.FormatBool(u.IsSuperuser)).Err(); err != nil {
		return fmt.Errorf("failed to save into redis %w", err)
	}

	return nil
}

func (s *Store) Load(ctx context.Context, username string) (User, error) {
	val, err := s.Client.HGetAll(ctx, Key(username)).Result()
	if err != nil {
		return User{}, fmt.Errorf("failed to load from redis %w", err)
	}

	if len(val) == 0 {
		return User{}, ErrUserNotExists
	}

	password, ok := val[PasswordHKey]
	if !ok {
		return User{}, ErrInvalidHResult
	}

	siss, ok := val[IsSuperuserHKey]
	if !ok {
		return User{}, ErrInvalidHResult
	}

	iss, err := strconv.ParseBool(siss)
	if err != nil {
		return User{}, fmt.Errorf("invalid is_superuser %w", err)
	}

	return User{
		Password:    password,
		IsSuperuser: iss,
		Username:    username,
	}, nil
}

func (s *Store) LoadAll(ctx context.Context) ([]User, error) {
	keys, err := s.Client.Keys(ctx, "mqtt_user:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys from redis %w", err)
	}

	users := make([]User, 0)

	for _, key := range keys {
		user, err := s.Load(ctx, strings.TrimPrefix(key, "mqtt_user:"))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user from redis %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}
