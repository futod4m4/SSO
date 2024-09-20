package validations

import (
	ssov1 "github.com/futod4m4/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login Handler validations

func ValidateLogin(req *ssov1.LoginRequest, validate *validator.Validate) error {
	if err := validateLoginEmail(req.GetEmail(), validate); err != nil {
		return err
	}

	if err := validateLoginPassword(req.GetPassword(), validate); err != nil {
		return err
	}

	if err := validateLoginAppId(req.GetAppId(), validate); err != nil {
		return err
	}

	return nil
}

// validateLoginEmail validates if email is correct
func validateLoginEmail(email string, validate *validator.Validate) error {

	if err := validate.Var(email, "email"); err != nil {
		return status.Error(codes.InvalidArgument, "incorrect email")
	}

	if err := validate.Var(email, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}

// validateLoginPassword validates if password is not empty and not too easy
func validateLoginPassword(password string, validate *validator.Validate) error {

	if err := validate.Var(password, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// validateLoginAppId validates is app_id is not 0
func validateLoginAppId(appId int32, validate *validator.Validate) error {

	if err := validate.Var(appId, "ne=0"); err != nil {
		return status.Error(codes.InvalidArgument, "incorrect app_id")
	}

	return nil
}

// Register Handler validations

// ValidateRegister validates register Handler
func ValidateRegister(req *ssov1.RegisterRequest, validate *validator.Validate) error {
	if err := validateRegisterEmail(req.GetEmail(), validate); err != nil {
		return err
	}

	if err := validateRegisterPassword(req.GetPassword(), validate); err != nil {
		return err
	}

	return nil
}

// validateRegisterEmail validates if email is correct
func validateRegisterEmail(email string, validate *validator.Validate) error {

	if err := validate.Var(email, "email"); err != nil {
		return status.Error(codes.InvalidArgument, "incorrect email")
	}

	if err := validate.Var(email, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}

// validateRegisterPassword validates if password is not empty and not too easy
func validateRegisterPassword(password string, validate *validator.Validate) error {
	if err := validate.Var(password, "lt=6"); err != nil {
		return status.Error(codes.InvalidArgument, "password is too short")
	}

	//Password should contain uppercase and lowercase letters and symbols
	if err := validate.Var(password, "excludesall=!@#?,lowercase|uppercase"); err != nil {
		return status.Error(codes.InvalidArgument, "password is too easy")
	}

	if err := validate.Var(password, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// IsAdmin Handler validations

// ValidateIsAdmin
func ValidateIsAdmin(userId int64, validate *validator.Validate) error {

	if err := validate.Var(userId, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
