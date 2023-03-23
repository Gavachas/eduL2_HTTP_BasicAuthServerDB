package repository

type Repository struct {
	drv Driver
}

func NewRepository(drv Driver) *Repository {
	return &Repository{drv}
}

func (r *Repository) InsertIncidetRep(name string, author int) (int, error) {
	return r.drv.InsertIncidet(name, author)
}
func (r *Repository) GetIncidentRep(id int) (*Incident, error) {
	return r.drv.GetIncident(id)
}
func (r *Repository) InsertUserRep(name, pass string) (int, error) {
	return r.drv.InsertUser(name, pass)
}
func (r *Repository) GetUserByNameRep(name string) (*User, error) {
	return r.drv.GetUserByName(name)
}
func (r *Repository) GetUserFirstRep() (int, error) {
	return r.drv.GetUserFirst()
}
func (r *Repository) InsertRulesRep(name string, user int) (int, error) {
	return r.drv.InsertRules(name, user)
}
func (r *Repository) GetUserRulesRep(id int) (*UserRules, error) {
	return r.drv.GetUserRules(id)
}
