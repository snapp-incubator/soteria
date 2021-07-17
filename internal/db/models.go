package db

import (
	"context"
	"errors"
	"time"
)

// ErrDb is returned when there was an error in database operations
var ErrDb = errors.New("database error")

type ModelHandler interface {
	Save(ctx context.Context, model Model) error
	Delete(ctx context.Context, modelName, pk string) error
	Get(ctx context.Context, modelName, pk string, model Model) error
	Update(ctx context.Context, model Model) error
}

type Model interface {
	GetMetadata() MetaData
	GetPrimaryKey() string
}

type MetaData struct {
	ModelName    string    `json:"model_name"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified"`
}
