package sqlDB

import "eduL2_HTTP_BasicAuthServerDB/internal/models"

type Driver interface {
	InsertIncidet(name string, author int) (int, error)
	GetIncident(id int) (*models.Incident, error)
	InsertUser(name, pass string) (int, error)
	GetUserByName(name string) (*models.User, error)
	GetUserFirst() (int, error)
	InsertRules(name string, user int) (int, error)
	GetUserRules(id int) (*models.UserRules, error)
	InitBase()
	InitBaseMain() error
}
