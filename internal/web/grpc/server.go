package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.snapp.ir/dispatching/soteria/v3/internal"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/web/grpc/contracts"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	contracts.AuthContract
}

func (s *Server) Auth(ctx context.Context, in *contracts.AuthContract) (*contracts.ServiceResponse, error) {
	start := time.Now()

	username := in.GetUsername()
	password := in.GetPassword()
	endpoint := in.GetEndpoint()
	ip := in.GetIPAddress()
	zap.L().Debug("grpc auth call",
		zap.String("username", username),
		zap.String("endpoint", endpoint),
		zap.String("ip", ip),
	)

	var ok bool
	var err error
	if len(password) > 0 {
		ok, err = app.GetInstance().Authenticator.EndPointBasicAuth(ctx, username, password, endpoint)
	} else if len(ip) > 0 {
		ok, err = app.GetInstance().Authenticator.EndpointIPAuth(ctx, username, ip, endpoint)
	} else {
		ok = false
		err = fmt.Errorf("both password and ip address are empty")
	}

	if !ok {
		if errors.Is(err, db.ErrDb) {
			app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Auth, http.StatusInternalServerError)
			app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Auth, internal.Failure, "database error happened")
			app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Auth, float64(time.Since(start).Nanoseconds()))

			zap.L().Error("grpc auth returned", zap.Int("code", http.StatusInternalServerError), zap.Error(err))
			return &contracts.ServiceResponse{Code: http.StatusInternalServerError}, fmt.Errorf("internal server error")
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Auth, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Auth, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Auth, float64(time.Since(start).Nanoseconds()))

		zap.L().Error("grpc auth returned", zap.Int("code", http.StatusUnauthorized), zap.Error(err))
		return &contracts.ServiceResponse{Code: http.StatusUnauthorized}, fmt.Errorf("request is unauthorized")
	}

	zap.L().Info("grpc auth returned", zap.Int("code", http.StatusOK))

	app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Auth, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Auth, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Auth, float64(time.Since(start).Nanoseconds()))

	return &contracts.ServiceResponse{Code: http.StatusOK}, nil
}

func (s *Server) GetToken(ctx context.Context, in *contracts.GetTokenContract) (*contracts.GetTokenResponse, error) {
	start := time.Now()

	grantType := in.GetGrantType()
	clientID := in.GetClientID()
	clientSecret := in.GetClientSecret()
	zap.L().Debug("grpc token call",
		zap.String("grant_type", grantType),
		zap.String("client_id", clientID),
		zap.String("client_secret", clientSecret),
	)

	tokenString, err := app.GetInstance().Authenticator.Token(ctx, acl.AccessType(grantType), clientID, clientSecret)
	if err != nil {
		if errors.Is(err, db.ErrDb) {
			app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Token, http.StatusInternalServerError)
			app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, "database error happened")
			app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, clientID)
			app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Token, float64(time.Since(start).Nanoseconds()))

			return &contracts.GetTokenResponse{
				Code:  http.StatusInternalServerError,
				Token: "",
			}, fmt.Errorf("internal server error")
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Token, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, clientID)
		app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Token, float64(time.Since(start).Nanoseconds()))

		return &contracts.GetTokenResponse{
			Code:  http.StatusUnauthorized,
			Token: "",
		}, fmt.Errorf("request is unauthorized")
	}

	zap.L().
		Info("token request accepted",
			zap.String("grant_type", grantType),
			zap.String("client_id", clientID),
			zap.String("client_secret", clientSecret),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Token, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Success, "token request accepted")
	app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Success, clientID)
	app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Token, float64(time.Since(start).Nanoseconds()))

	return &contracts.GetTokenResponse{
		Code:  200,
		Token: tokenString,
	}, nil
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	contracts.RegisterAuthServer(s, &Server{})
	return s
}
