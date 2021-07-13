package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Create is the body payload structure of create emq endpoint.
type Create struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Duration int64  `json:"duration"`
}

func (r Create) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Duration, validation.Required),
	); err != nil {
		return fmt.Errorf("create request validation failed: %w", err)
	}

	return nil
}
