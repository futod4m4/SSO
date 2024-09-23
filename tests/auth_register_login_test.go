package tests

import (
	"SSO/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	ssov1 "github.com/futod4m4/protos/gen/go/sso"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppId = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()
	userName := gofakeit.Username()
	location := gofakeit.Country()
	dateOfBirth := gofakeit.Date().Format("2006-01-02")
	sex := ""

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:       email,
		Password:    password,
		Username:    userName,
		Sex:         sex,
		Location:    location,
		DateOfBirth: dateOfBirth,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegister_DoubleRegister(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()
	userName := gofakeit.Username()
	location := gofakeit.Country()
	dateOfBirth := gofakeit.Date().Format("2006-01-02")
	sex := "Male"

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:       email,
		Password:    password,
		Username:    userName,
		Sex:         sex,
		Location:    location,
		DateOfBirth: dateOfBirth,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg2, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:       email,
		Password:    password,
		Username:    userName,
		Sex:         sex,
		Location:    location,
		DateOfBirth: dateOfBirth,
	})
	require.Error(t, err)
	require.Empty(t, respReg2.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
