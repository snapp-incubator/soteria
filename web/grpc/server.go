package grpc

import (
	"context"
	"errors"
	"fmt"
	"gitlab.snapp.ir/dispatching/soteria/internal"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/web/grpc/contracts"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"time"
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
		zap.String("password", password),
		zap.String("endpoint", endpoint),
		zap.String("ip", ip),
	)

	var ok bool
	var err error
	if len(password) > 0 {
		ok, err = app.GetInstance().Authenticator.EndPointBasicAuth(username, password, endpoint)
	} else if len(ip) > 0 {
		ok, err = app.GetInstance().Authenticator.EndpointIPAuth(username, ip, endpoint)
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
			return &contracts.ServiceResponse{Code: http.StatusInternalServerError}, err
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Auth, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Auth, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Auth, float64(time.Since(start).Nanoseconds()))

		zap.L().Error("grpc auth returned", zap.Int("code", http.StatusUnauthorized), zap.Error(err))
		return &contracts.ServiceResponse{Code: http.StatusUnauthorized}, err
	}

	zap.L().Info("grpc auth returned", zap.Int("code", http.StatusOK))

	app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Auth, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Auth, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Auth, float64(time.Since(start).Nanoseconds()))

	return &contracts.ServiceResponse{Code: http.StatusOK}, err
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

	tokenString, err := app.GetInstance().Authenticator.Token(acl.AccessType(grantType), clientID, clientSecret)
	if err != nil {
		if errors.Is(err, db.ErrDb) {
			app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Token, http.StatusInternalServerError)
			app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, "database error happened")
			app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, clientID)
			app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Token, float64(time.Since(start).Nanoseconds()))

			return &contracts.GetTokenResponse{
				Code:  http.StatusInternalServerError,
				Token: "",
			}, err
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.GrpcApi, internal.Soteria, internal.Token, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveStatus(internal.GrpcApi, internal.Soteria, internal.Token, internal.Failure, clientID)
		app.GetInstance().Metrics.ObserveResponseTime(internal.GrpcApi, internal.Soteria, internal.Token, float64(time.Since(start).Nanoseconds()))

		return &contracts.GetTokenResponse{
			Code:  http.StatusUnauthorized,
			Token: "",
		}, err
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

func GRPCServer() *grpc.Server {
	s := grpc.NewServer()
	contracts.RegisterAuthServer(s, &Server{})
	return s
}
