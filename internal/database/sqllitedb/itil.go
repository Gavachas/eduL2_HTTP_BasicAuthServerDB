package sqllitedb

import (
	"database/sql"
	"eduL2_HTTP_BasicAuthServerDB/internal/models"
	"errors"
	"log"
)

type ItilModel struct {
	db *sql.DB
}

func NewDriver(db *sql.DB) *ItilModel {
	return &ItilModel{db: db}
}
func (m *ItilModel) InsertIncidet(name string, author int) (int, error) {
	stmt := "INSERT INTO incidents (name, author)  VALUES(?, ?)"
	result, err := m.db.Exec(stmt, name, author)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (m *ItilModel) GetIncident(id int) (*models.Incident, error) {
	stmt := `SELECT inc.id, inc.name, u.id  FROM incidents AS  inc
	LEFT JOIN users    AS  u on u.id = inc.author 
	WHERE inc.id  = ?`
	row := m.db.QueryRow(stmt, id)
	inc := &models.Incident{}
	err := row.Scan(&inc.Id, &inc.Name, &inc.Author)
	//log.Fatal("err   ", err, inc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return inc, nil
}
func (m *ItilModel) InsertUser(name, pass string) (int, error) {
	stmt := "INSERT INTO users (name, pass)  VALUES(?, ?)"
	result, err := m.db.Exec(stmt, name, pass)
	if err != nil {

		return 0, err

	}
	id, err := result.LastInsertId()
	if err != nil {

		return 0, err

	}
	return int(id), nil
}
func (m *ItilModel) GetUserByName(name string) (*models.User, error) {
	stmt := `SELECT id, name, pass  FROM users   
	WHERE name  = ?`
	row := m.db.QueryRow(stmt, name)
	user := &models.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Pass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return user, nil

}
func (m *ItilModel) GetUserFirst() (int, error) {
	stmt := "SELECT COUNT(*) FROM users "
	row := m.db.QueryRow(stmt)
	var userID int
	userID = -1
	err := row.Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return userID, models.ErrNoRecord
		} else {

			return userID, err
		}
	}
	return userID, nil

}
func (m *ItilModel) InsertRules(name string, user int) (int, error) {
	stmt := "INSERT INTO rules (name, userID)  VALUES(?, ?)"
	result, err := m.db.Exec(stmt, name, user)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (m *ItilModel) GetUserRules(id int) (*models.UserRules, error) {
	stmt := `SELECT id, name  FROM rules   
	WHERE userID  = ?`
	row := m.db.QueryRow(stmt, id)
	userRules := &models.UserRules{}
	err := row.Scan(&userRules.Id, &userRules.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	userRules.User = id
	return userRules, nil

}
func (m *ItilModel) InitBase() {
	// We store bcrypt-ed passwords for each user. The actual passwords are "1234"
	// for "joe" and "strongerpassword9902" for "mary", but these should not be
	// stored anywhere.
	var usersPasswords = map[string][]byte{
		"joe":  []byte("$2a$12$aMfFQpGSiPiYkekov7LOsu63pZFaWzmlfm1T8lvG6JFj2Bh4SZPWS"),
		"mary": []byte("$2a$12$l398tX477zeEBP6Se0mAv.ZLR8.LZZehuDgbtw2yoQeMjIyCNCsRW"),
	}
	var usersRules = map[string]string{
		"joe":  "admin",
		"mary": "user",
	}
	for user, pass := range usersPasswords {
		m.InsertUser(user, string(pass))
		log.Println("Init db  ", user)
	}

	for user, rule := range usersRules {
		uss, err := m.GetUserByName(user)
		if err != nil {
			continue
		}
		m.InsertRules(rule, uss.Id)
		log.Println("Init db  ", rule)
	}
}
func (m *ItilModel) InitBaseMain() error {
	var err error
	_, err = m.db.Exec("CREATE DATABASE Itil")
	if err != nil {
		return err
	}

	_, err = m.db.Exec("USE Itil")
	if err != nil {
		return err
	}
	return nil
}
