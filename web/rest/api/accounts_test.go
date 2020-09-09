package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	redisModel "gitlab.snapp.ir/dispatching/soteria/internal/db/redis"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	app.GetInstance().SetAccountsService(&accounts.Service{
		Handler: redisModel.ModelHandler{
			Client: client,
		},
	})
}

func TestCreateAccount(t *testing.T) {
	router := setupRouter()

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"username":"user","password":"123","user_type":"driver"}`)
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		expectedResponse := Response{
			Code:    errors.SuccessfulOperation,
			Message: fmt.Sprintf("%s: []", errors.SuccessfulOperation.Message()),
			Data:    nil,
		}

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("testing without content type header", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"username":"user","password":"123","user_type":"driver"}`)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(payload))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.BadRequestPayload, actualResponse.Code)
	})
}

func TestReadAccount(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/accounts/user", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		u := actualResponse.Data.(map[string]interface{})
		assert.Equal(t, "user", u["username"])
		assert.Equal(t, "passenger", u["type"])
	})
}

func TestUpdateAccount(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"new_password":"password2","ips":["127.0.0.1"],"secret":"12345678","type":"EMQUser","token_expiration":1000000}`)
		req, _ := http.NewRequest(http.MethodPut, "/accounts/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		u, err := app.GetInstance().AccountsService.Info("user", "password2")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(u.IPs))
		assert.Equal(t, "12345678", u.Secret)
		assert.Equal(t, user.EMQUser, u.Type)
		assert.Equal(t, time.Duration(1000000), u.TokenExpirationDuration)
	})
}

func TestDeleteAccount(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/accounts/user", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		_, err = app.GetInstance().AccountsService.Info("user", "password2")
		assert.NotNil(t, err)
	})
}

func TestUpdatePublicKey(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")

	t.Run("testing with an invalid public key", func(t *testing.T) {
		w := httptest.NewRecorder()

		file, err := ioutil.ReadFile("../../../test/token.invalid.sample")
		assert.NoError(t, err)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("public_key", "token.invalid.sample")
		assert.NoError(t, err)
		part.Write(file)
		err = writer.Close()
		assert.NoError(t, err)
		req, _ := http.NewRequest(http.MethodPost, "/accounts/user/key", body)

		req.SetBasicAuth("user", "password")
		req.Header.Set("Content-Type", writer.FormDataContentType())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.PublicKeyParseFailure, actualResponse.Code)
	})

	t.Run("testing with an invalid form key", func(t *testing.T) {
		w := httptest.NewRecorder()

		file, err := ioutil.ReadFile("../../../test/0.pem")
		assert.NoError(t, err)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("key", "0.pem")
		assert.NoError(t, err)
		part.Write(file)
		err = writer.Close()
		assert.NoError(t, err)
		req, _ := http.NewRequest(http.MethodPost, "/accounts/user/key", body)

		req.SetBasicAuth("user", "password")
		req.Header.Set("Content-Type", writer.FormDataContentType())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.PublicKeyReadFormFailure, actualResponse.Code)
	})

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		file, err := ioutil.ReadFile("../../../test/0.pem")
		assert.NoError(t, err)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("public_key", "0.pem")
		assert.NoError(t, err)
		part.Write(file)
		err = writer.Close()
		assert.NoError(t, err)
		req, _ := http.NewRequest(http.MethodPost, "/accounts/user/key", body)

		req.SetBasicAuth("user", "password")
		req.Header.Set("Content-Type", writer.FormDataContentType())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		u, err := app.GetInstance().AccountsService.Info("user", "password")
		assert.Nil(t, err)
		assert.NotNil(t, u.PublicKey)
	})
}

func TestCreateAccountRule(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")

	t.Run("testing with invalid rule info", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"endpoint":"/notification","topic":"","access_type":"2"}`)
		req, _ := http.NewRequest(http.MethodPost, "/accounts/user/rules", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.InvalidRule, actualResponse.Code)

		u, err := app.GetInstance().AccountsService.Info("user", "password")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(u.Rules))
	})

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"endpoint":"/notification","topic":"","access_type":""}`)
		req, _ := http.NewRequest(http.MethodPost, "/accounts/user/rules", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		u, err := app.GetInstance().AccountsService.Info("user", "password")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(u.Rules))
		assert.Equal(t, "/notification", u.Rules[0].Endpoint)
		assert.Equal(t, topics.Type(""), u.Rules[0].Topic)
		assert.Equal(t, acl.AccessType(""), u.Rules[0].AccessType)
	})
}

func TestReadAccountRule(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")
	createdRule, _ := app.GetInstance().AccountsService.CreateRule("user", "/notification", "", "")

	t.Run("testing with invalid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/accounts/user/rules/invalid-uuid", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.InvalidRuleUUID, actualResponse.Code)
	})

	t.Run("testing with undefined rule", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/accounts/user/rules/b33a0b78-c8a6-4719-a222-9a3883cc4b7c", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.RuleNotFound, actualResponse.Code)
	})

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/user/rules/%s", createdRule.UUID), nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		returnedRule := actualResponse.Data.(map[string]interface{})
		assert.Equal(t, createdRule.UUID.String(), returnedRule["uuid"])
		assert.Equal(t, "/notification", returnedRule["endpoint"])
		assert.Equal(t, "", returnedRule["topic"])
		assert.Equal(t, "", returnedRule["topic"])
	})
}

func TestUpdateAccountRule(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")
	createdRule, _ := app.GetInstance().AccountsService.CreateRule("user", "/notification", "", "")

	t.Run("testing with no payload", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPut, "/accounts/user/rules/b33a0b78-c8a6-4719-a222-9a3883cc4b7c", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.BadRequestPayload, actualResponse.Code)
	})

	t.Run("testing with invalid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"endpoint":"/notification","topic":"","access_type":""}`)
		req, _ := http.NewRequest(http.MethodPut, "/accounts/user/rules/invalid-uuid", bytes.NewBuffer(payload))
		req.SetBasicAuth("user", "password")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.InvalidRuleUUID, actualResponse.Code)
	})

	t.Run("testing with undefined rule", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"endpoint":"/notification","topic":"","access_type":""}`)
		req, _ := http.NewRequest(http.MethodPut, "/accounts/user/rules/b33a0b78-c8a6-4719-a222-9a3883cc4b7c", bytes.NewBuffer(payload))
		req.SetBasicAuth("user", "password")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.RuleNotFound, actualResponse.Code)
	})

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"endpoint":"","topic":"cab_event","access_type":"2"}`)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/accounts/user/rules/%s", createdRule.UUID), bytes.NewBuffer(payload))
		req.SetBasicAuth("user", "password")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		u, err := app.GetInstance().AccountsService.Info("user", "password")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(u.Rules))
		assert.Equal(t, createdRule.UUID, u.Rules[0].UUID)
		assert.Equal(t, "", u.Rules[0].Endpoint)
		assert.Equal(t, topics.CabEvent, u.Rules[0].Topic)
		assert.Equal(t, acl.Pub, string(u.Rules[0].AccessType))
	})
}

func TestDeleteAccountRule(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")
	createdRule, _ := app.GetInstance().AccountsService.CreateRule("user", "/notification", "", "")

	t.Run("testing with invalid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, "/accounts/user/rules/invalid-uuid", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.InvalidRuleUUID, actualResponse.Code)
	})

	t.Run("testing with undefined rule", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, "/accounts/user/rules/b33a0b78-c8a6-4719-a222-9a3883cc4b7c", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.RuleNotFound, actualResponse.Code)
	})

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/accounts/user/rules/%s", createdRule.UUID), nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		resBody, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var actualResponse Response
		err = json.Unmarshal(resBody, &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, errors.SuccessfulOperation, actualResponse.Code)

		u, err := app.GetInstance().AccountsService.Info("user", "password")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(u.Rules))
	})
}
