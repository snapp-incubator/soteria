package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountsBasicAuth(t *testing.T) {
	router := gin.Default()
	router.Use(accountsBasicAuth())
	router.GET("/:username", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "done")
		return
	})

	_ = app.GetInstance().AccountsService.SignUp(context.Background(), "user", "password", "passenger")

	t.Run("testing successful authentication", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/user", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("testing without a Authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("testing with an invalid auth", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/user", nil)
		req.Header.Set("Authorization", "Basic invalid")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("testing different username in path", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/user2", nil)
		req.SetBasicAuth("user", "password")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("testing with wrong password", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/user", nil)
		req.SetBasicAuth("user", "password2")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
