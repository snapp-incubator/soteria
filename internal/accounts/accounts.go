package accounts

import "gitlab.snapp.ir/dispatching/soteria/v3/internal/db"

// Service is responsible for all things related to account handling.
type Service struct {
	Handler db.ModelHandler
}
