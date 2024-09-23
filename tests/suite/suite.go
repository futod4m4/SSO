package suite

import (
	"SSO/internal/config"
	"context"
	ssov1 "github.com/futod4m4/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T                  // Потребуется для вызова методов *testing.T
	Cfg        *config.Config   // Конфигурация приложения
	AuthClient ssov1.AuthClient // Клиент для взаимодействия с gRPC-сервером
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAdress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAdress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
