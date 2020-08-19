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
	"gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
		Handler: redisModel.RedisModelHandler{
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

		assert.Equal(t, errors.BadRequestPayload, string(actualResponse.Code))
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

		user := actualResponse.Data.(map[string]interface{})
		assert.Equal(t, "user", user["username"])
		assert.Equal(t, "passenger", user["type"])
	})
}

func TestUpdateAccount(t *testing.T) {
	router := setupRouter()

	_ = app.GetInstance().AccountsService.SignUp("user", "password", "passenger")

	t.Run("testing successful request", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload := []byte(`{"new_password":"password2"}`)
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

		_, err = app.GetInstance().AccountsService.Info("user", "password2")
		assert.Nil(t, err)
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
