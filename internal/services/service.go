package services

import (
	"eduL2_HTTP_BasicAuthServerDB/internal/repository"
	"log"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	Rep      *repository.Repository
}

func NewService(errorLog, infoLog *log.Logger, rep *repository.Repository) *application {
	return &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		Rep:      rep,
	}
}
