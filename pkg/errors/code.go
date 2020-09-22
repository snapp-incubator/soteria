package errors

import "net/http"

// Code is the type of the error returned to user in api
type Code string

const (
	SuccessfulOperation           Code = "successful_operation"
	BadRequestPayload             Code = "bad_request_payload"
	SignUpUserFailure             Code = "sign_up_user_failure"
	WrongUsernameOrPassword       Code = "wrong_username_or_password"
	DatabaseSaveFailure           Code = "database_save_failure"
	DatabaseGetFailure            Code = "database_get_failure"
	DatabaseUpdateFailure         Code = "database_update_failure"
	DatabaseDeleteFailure         Code = "database_delete_failure"
	PasswordHashGenerationFailure Code = "password_hash_generation_failure"
	UsernameMismatch              Code = "username_mismatch"
	IPMisMatch                    Code = "ip_mismatch"
	PublicKeyReadFormFailure      Code = "public_key_read_form_failure"
	PublicKeyOpenFailure          Code = "public_key_open_failure"
	PublicKeyReadFileFailure      Code = "public_key_read_file_failure"
	PublicKeyParseFailure         Code = "public_key_parse_failure"
	InvalidRuleUUID               Code = "invalid_rule_uuid"
	RuleNotFound                  Code = "rule_not_found"
	InvalidRule                   Code = "invalid_rule"
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
	case IPMisMatch:
		return "ip is not authorized"
	case PublicKeyReadFormFailure:
		return "cloud not read public key file from form"
	case PublicKeyOpenFailure:
		return "could not open public key file"
	case PublicKeyReadFileFailure:
		return "could not read from opened public key file"
	case PublicKeyParseFailure:
		return "could not parse public key"
	case InvalidRuleUUID:
		return "provided rule UUID is not valid"
	case RuleNotFound:
		return "account has no rule with provided UUID"
	case InvalidRule:
		return "provided rule is not valid"
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
	case IPMisMatch:
		return http.StatusUnauthorized
	case PublicKeyReadFormFailure:
		return http.StatusBadRequest
	case PublicKeyOpenFailure:
		return http.StatusBadRequest
	case PublicKeyReadFileFailure:
		return http.StatusInternalServerError
	case PublicKeyParseFailure:
		return http.StatusBadRequest
	case InvalidRuleUUID:
		return http.StatusBadRequest
	case RuleNotFound:
		return http.StatusNotFound
	case InvalidRule:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
