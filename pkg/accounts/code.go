package accounts

import "net/http"

// Code is the type of the error returned to user in REST API
type Code string

const (
	SuccessfulOperation           Code = "successful_operation"
	BadRequestPayload                  = "bad_request_payload"
	SignUpUserFailure                  = "sign_up_user_failure"
	WrongUsernameOrPassword            = "wrong_username_or_password"
	DatabaseSaveFailure                = "database_save_failure"
	DatabaseGetFailure                 = "database_get_failure"
	DatabaseUpdateFailure              = "database_update_failure"
	DatabaseDeleteFailure              = "database_delete_failure"
	PasswordHashGenerationFailure      = "password_hash_generation_failure"
	UsernameMismatch                   = "username_mismatch"
)

// Message will return detailed error information
func (e Code) Message() string {
	switch e {
	case SuccessfulOperation:
		return "operation done successfully"
	case BadRequestPayload:
		return "request payload is not valid"
	case SignUpUserFailure:
		return "could not save user"
	case WrongUsernameOrPassword:
		return "username or password is not correct"
	case DatabaseSaveFailure:
		return "could not save info to database"
	case DatabaseGetFailure:
		return "could not get info from database"
	case DatabaseUpdateFailure:
		return "could update info in database"
	case DatabaseDeleteFailure:
		return "could delete info from database"
	case PasswordHashGenerationFailure:
		return "could not generate hash from password"
	case UsernameMismatch:
		return "provided username does not match with username in path"
	default:
		return ""
	}
}

// HttpStatusCode returns the perfect HTTP Status Code for the error
func (e Code) HttpStatusCode() int {
	switch e {
	case SuccessfulOperation:
		return http.StatusOK
	case BadRequestPayload:
		return http.StatusBadRequest
	case SignUpUserFailure:
		return http.StatusInternalServerError
	case WrongUsernameOrPassword:
		return http.StatusUnauthorized
	case DatabaseSaveFailure:
		return http.StatusInternalServerError
	case DatabaseGetFailure:
		return http.StatusInternalServerError
	case DatabaseUpdateFailure:
		return http.StatusInternalServerError
	case DatabaseDeleteFailure:
		return http.StatusInternalServerError
	case PasswordHashGenerationFailure:
		return http.StatusInternalServerError
	case UsernameMismatch:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
