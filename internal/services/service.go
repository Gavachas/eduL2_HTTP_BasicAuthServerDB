package services

import (
	"eduL2_HTTP_BasicAuthServerDB/internal/repository"
	"eduL2_HTTP_BasicAuthServerDB/pkg/sessiontoken"

	"go.uber.org/zap"
)

type application struct {
	logger  *zap.Logger
	Rep     *repository.Repository
	Session *sessiontoken.Session
}

func NewService(elogger *zap.Logger, rep *repository.Repository, session *sessiontoken.Session) *application {
	return &application{
		logger:  elogger,
		Rep:     rep,
		Session: session,
	}
}
