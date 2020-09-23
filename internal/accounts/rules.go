package accounts

import (
	"github.com/google/uuid"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"time"
)

// CreateRule will add a new rule with given information to a user
func (s Service) CreateRule(username, endpoint string, topicPattern topics.Type, accessType acl.AccessType) (*user.Rule, *errors.Error) {
	if err := validateRuleInfo(endpoint, topicPattern, accessType); err != nil {
		return nil, err
	}

	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return nil, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	rule := &user.Rule{
		UUID:       uuid.New(),
		Endpoint:   endpoint,
		Topic:      topicPattern,
		AccessType: accessType,
	}

	u.Rules = append(u.Rules, *rule)

	u.MetaData.DateModified = time.Now()

	if err := s.Handler.Update(u); err != nil {
		return nil, errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return rule, nil
}

// GetRule returns a user rule based on given username, password and rule UUID
func (s Service) GetRule(username string, id uuid.UUID) (*user.Rule, *errors.Error) {
	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return nil, errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	for _, r := range u.Rules {
		if r.UUID == id {
			return &r, nil
		}
	}

	return nil, errors.CreateError(errors.RuleNotFound, "")
}

// UpdateRule updates an account's rule based on given username, UUID and information
func (s Service) UpdateRule(username string, id uuid.UUID, endpoint string, topicPattern topics.Type, accessType acl.AccessType) *errors.Error {
	if err := validateRuleInfo(endpoint, topicPattern, accessType); err != nil {
		return err
	}

	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	var ruleIndex *int
	for i, r := range u.Rules {
		if r.UUID == id {
			ruleIndex = &i
		}
	}
	if ruleIndex == nil {
		return errors.CreateError(errors.RuleNotFound, "")
	}

	u.Rules[*ruleIndex].Endpoint = endpoint
	u.Rules[*ruleIndex].Topic = topicPattern
	u.Rules[*ruleIndex].AccessType = accessType

	u.MetaData.DateModified = time.Now()

	if err := s.Handler.Update(u); err != nil {
		return errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return nil
}

// DeleteRule deletes an account's rule based on given username and UUID
func (s Service) DeleteRule(username string, id uuid.UUID) *errors.Error {
	var u user.User
	if err := s.Handler.Get("user", username, &u); err != nil {
		return errors.CreateError(errors.DatabaseGetFailure, err.Error())
	}

	var ruleIndex *int
	for i, r := range u.Rules {
		if r.UUID == id {
			ruleIndex = &i
		}
	}
	if ruleIndex == nil {
		return errors.CreateError(errors.RuleNotFound, "")
	}

	u.Rules = append(u.Rules[:*ruleIndex], u.Rules[*ruleIndex+1:]...)

	u.MetaData.DateModified = time.Now()

	if err := s.Handler.Update(u); err != nil {
		return errors.CreateError(errors.DatabaseUpdateFailure, err.Error())
	}

	return nil
}

// validateRuleInfo will validate a rule info
func validateRuleInfo(endpoint string, topicPattern topics.Type, accessType acl.AccessType) *errors.Error {
	if endpoint == "" && topicPattern == "" && accessType == "" {
		return errors.CreateError(errors.InvalidRule, "all rule information is empty")
	}

	if endpoint != "" && (topicPattern != "" || accessType != "") {
		return errors.CreateError(errors.InvalidRule, "when endpoint is provided topic pattern and access type should be empty")
	}

	if (topicPattern != "" && accessType == "") || (topicPattern == "" && accessType != "") {
		return errors.CreateError(errors.InvalidRule, "both topic pattern and access type should be present")
	}

	return nil
}