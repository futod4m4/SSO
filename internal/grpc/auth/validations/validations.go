package validations

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login Handler validations

// ValidateLoginEmail validates if email is correct
func ValidateLoginEmail(email string, validate *validator.Validate) error {

	if err := validate.Var(email, "email"); err != nil {
		return status.Error(codes.InvalidArgument, "incorrect email")
	}

	if err := validate.Var(email, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}

// ValidateLoginPassword validates if password is not empty and not too easy
func ValidateLoginPassword(password string, validate *validator.Validate) error {

	if err := validate.Var(password, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// ValidateLoginAppId validates is app_id is not 0
func ValidateLoginAppId(appId int32, validate *validator.Validate) error {

	if err := validate.Var(appId, "ne=0"); err != nil {
		return status.Error(codes.InvalidArgument, "incorrect app_id")
	}

	return nil
}

// Register Handler validations

// ValidateRegisterEmail validates if email is correct
func ValidateRegisterEmail(email string, validate *validator.Validate) error {

	if err := validate.Var(email, "email"); err != nil {
		return status.Error(codes.InvalidArgument, "incorrect email")
	}

	if err := validate.Var(email, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}

// ValidateRegisterPassword validates if password is not empty and not too easy
func ValidateRegisterPassword(password string, validate *validator.Validate) error {
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
