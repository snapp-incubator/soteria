package grpc

import (
	"context"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/web/grpc/contracts"
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
	statusCode := http.StatusUnauthorized
	ok := false
	var err error
	if len(password) > 0 {
		ok, err = app.GetInstance().Authenticator.EndPointBasicAuth(username, password, endpoint)
	} else if len(ip) > 0 {
		ok, err = app.GetInstance().Authenticator.EndpointIPAuth(username, ip, endpoint)
	}
	if ok {
		statusCode = http.StatusOK
	}
	return &contracts.ServiceResponse{Code: int32(statusCode)}, err
}

func (s *Server) GetToken(ctx context.Context, in *contracts.GetTokenContract) (*contracts.GetTokenResponse, error) {
	tokenString, err := app.GetInstance().Authenticator.Token(in.GetGrantType(), in.GetClientID(), in.GetClientSecret())
	if err != nil {
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
