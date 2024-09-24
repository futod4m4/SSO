package auth

import (
	"SSO/internal/lib/validations"
	"SSO/internal/services/auth"
	"context"
	"errors"
	ssov1 "github.com/futod4m4/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appId int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
		username string,
		sex string,
		location string,
		dateOfBirth string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	IsUserExists(ctx context.Context, email string) (bool, error)
}

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if err := validations.ValidateLogin(req, validate); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "email or password is incorrect")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := validations.ValidateRegister(req, validate); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword(), req.GetUsername(), req.GetSex(), req.GetLocation(), req.GetDateOfBirth())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {

	if err := validations.ValidateIsAdmin(req.GetUserId(), validate); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) IsUserExists(
	ctx context.Context,
	req *ssov1.IsUserExistsRequest,
) (*ssov1.IsUserExistsResponse, error) {

	if err := validations.ValidateIsUserExists(req.GetEmail(), validate); err != nil {
		return nil, err
	}

	isExists, _ := s.auth.IsUserExists(ctx, req.GetEmail())

	return &ssov1.IsUserExistsResponse{
		IsExists: isExists,
	}, nil
}
