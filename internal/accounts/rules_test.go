package accounts

import (
	"github.com/stretchr/testify/assert"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"testing"
)

func TestValidateRuleInfo(t *testing.T) {
	t.Run("test with with invalid rule info", func(t *testing.T) {
		endpoint := ""
		topicPattern := topics.Type("")
		accessType := acl.AccessType("")

		err := validateRuleInfo(endpoint, topicPattern, accessType)
		assert.NotNil(t, err)
		assert.Equal(t, errors.InvalidRule, err.Code)

		endpoint = "/notification"
		topicPattern = topics.DriverLocation
		accessType = ""

		err = validateRuleInfo(endpoint, topicPattern, accessType)
		assert.NotNil(t, err)
		assert.Equal(t, errors.InvalidRule, err.Code)

		endpoint = "/notification"
		topicPattern = topics.DriverLocation
		accessType = ""

		err = validateRuleInfo(endpoint, topicPattern, accessType)
		assert.NotNil(t, err)
		assert.Equal(t, errors.InvalidRule, err.Code)

		endpoint = "/notification"
		topicPattern = topics.CabEvent
		accessType = acl.Pub

		err = validateRuleInfo(endpoint, topicPattern, accessType)
		assert.NotNil(t, err)
		assert.Equal(t, errors.InvalidRule, err.Code)

		endpoint = ""
		topicPattern = topics.DriverLocation
		accessType = ""

		err = validateRuleInfo(endpoint, topicPattern, accessType)
		assert.NotNil(t, err)
		assert.Equal(t, errors.InvalidRule, err.Code)

		endpoint = ""
		topicPattern = ""
		accessType = acl.Sub

		err = validateRuleInfo(endpoint, topicPattern, accessType)
		assert.NotNil(t, err)
		assert.Equal(t, errors.InvalidRule, err.Code)
	})

	t.Run("test with with valid rule info", func(t *testing.T) {
		endpoint := "/notification"
		topicPattern := topics.Type("")
		accessType := acl.Pub

		err := validateRuleInfo(endpoint, topicPattern, accessType)
		assert.Nil(t, err)

		endpoint = ""
		topicPattern = topics.DriverLocation
		accessType = acl.Pub

		err = validateRuleInfo(endpoint, topicPattern, accessType)
		assert.Nil(t, err)
	})
}
