package auth

import (
	"SSO/internal/grpc/auth/validations"
	"context"
	ssov1 "github.com/futod4m4/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

const (
	emptyValue = 0
)

var (
	valid = validator.New(validator.WithRequiredStructEnabled())
)

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if err := validations.ValidateLoginEmail(req.GetEmail(), valid); err != nil {
		return nil, err
	}

	if err := validations.ValidateLoginPassword(req.GetPassword(), valid); err != nil {
		return nil, err
	}

	if err := validations.ValidateLoginAppId(req.GetAppId(), valid); err != nil {
		return nil, err
	}

	return &ssov1.LoginResponse{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
