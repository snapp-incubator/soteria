package grpc

import (
	"context"
	"errors"
	"fmt"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/web/grpc/contracts"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

type Server struct {
	contracts.AuthContract
}

func (s *Server) Auth(ctx context.Context, in *contracts.AuthContract) (*contracts.ServiceResponse, error) {
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
			zap.L().Error("grpc auth returned", zap.Int("code", http.StatusInternalServerError), zap.Error(err))
			return &contracts.ServiceResponse{Code: int32(http.StatusInternalServerError)}, err
		}

		zap.L().Error("grpc auth returned", zap.Int("code", http.StatusUnauthorized), zap.Error(err))
		return &contracts.ServiceResponse{Code: int32(http.StatusUnauthorized)}, err
	}
	zap.L().Info("grpc auth returned", zap.Int("code", http.StatusOK))
	return &contracts.ServiceResponse{Code: int32(http.StatusOK)}, err
}

func (s *Server) GetToken(ctx context.Context, in *contracts.GetTokenContract) (*contracts.GetTokenResponse, error) {
	tokenString, err := app.GetInstance().Authenticator.Token(acl.AccessType(in.GetGrantType()), in.GetClientID(), in.GetClientSecret())
	if err != nil {
		if errors.Is(err, db.ErrDb) {
			return &contracts.GetTokenResponse{
				Code:  http.StatusInternalServerError,
				Token: "",
			}, err
		}

		return &contracts.GetTokenResponse{
			Code:  http.StatusUnauthorized,
			Token: "",
		}, err
	}
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
