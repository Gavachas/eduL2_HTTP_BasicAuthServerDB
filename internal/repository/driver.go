package repository

type Driver interface {
	StopConnect() error
	InsertIncidet(name string, author int) (int, error)
	GetIncident(id int) (*Incident, error)
	InsertUser(name, pass string) (int, error)
	GetUserByName(name string) (*User, error)
	GetUserFirst() (int, error)
	InsertRules(name string, user int) (int, error)
	GetUserRules(id int) (*UserRules, error)
}
