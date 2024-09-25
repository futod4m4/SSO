package app

import (
	grpcapp "SSO/internal/app/grpc"
	"SSO/internal/services/auth"
	"SSO/storage/postgresql"
	"fmt"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *postgresql.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgresql.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
		Storage: storage,
	}
}

func (a *App) Stop() error {
	var err error

	a.GRPCSrv.Stop()

	if err = a.Storage.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
