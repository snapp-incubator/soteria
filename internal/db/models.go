package db

import (
	"time"
)

type ModelHandler interface {
	Save(model Model) error
	Delete(modelName, pk string) error
	Get(modelName, pk string, v interface{}) error
	Update(model Model) error
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
