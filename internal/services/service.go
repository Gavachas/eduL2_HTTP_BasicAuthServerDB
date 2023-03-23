package services

import (
	"eduL2_HTTP_BasicAuthServerDB/internal/repository"

	"go.uber.org/zap"
)

type application struct {
	logger *zap.Logger
	Rep    *repository.Repository
}

func NewService(elogger *zap.Logger, rep *repository.Repository) *application {
	return &application{
		logger: elogger,
		Rep:    rep,
	}
}
