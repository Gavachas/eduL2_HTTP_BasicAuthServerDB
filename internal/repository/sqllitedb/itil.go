package sqllitedb

import (
	"context"
	"database/sql"
	"eduL2_HTTP_BasicAuthServerDB/internal/repository"
	"errors"
	"fmt"
	"log"

	sample_grpc "github.com/Gavachas/grpc_sample/grpc_s"
	"google.golang.org/grpc"
)

type ItilModel struct {
	db *sql.DB
}

func NewDriver(dbtype, dsn string) (*ItilModel, error) {
	im := &ItilModel{}
	err := im.connect(dbtype, dsn)
	if err != nil {
		return im, err
	}
	err = im.initBaseTable()
	return im, err
}

func (m *ItilModel) connect(dbtype, dsn string) error {
	var err error
	log.Println("err  ", dbtype, dsn)
	m.db, err = sql.Open(dbtype, dsn)
	if err != nil {
		return err
	}
	log.Println("err  ", err)
	if err = m.db.Ping(); err != nil {
		return err

	}
	log.Println("err  ", err)

	return err
}
func (m *ItilModel) StopConnect() error {
	var err error
	err = m.db.Close()
	return err
}
func (m *ItilModel) initBaseTable() error {
	var usId int
	var err error
	usId = -1
	if usId, err = m.GetUserFirst(); usId < 1 {
		if err != nil {
			log.Fatal(err)
			return err
		}
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
		log.Println("Try init db  ", usId)

		//log.Println("Init db  ")
	}
	return err
}
func (m *ItilModel) initBaseMain() error {
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
func (m *ItilModel) GetIncident(id int) (*repository.Incident, error) {
	stmt := `SELECT inc.id, inc.name, u.id  FROM incidents AS  inc
	LEFT JOIN users    AS  u on u.id = inc.author 
	WHERE inc.id  = ?`
	row := m.db.QueryRow(stmt, id)
	inc := &repository.Incident{}
	err := row.Scan(&inc.Id, &inc.Name, &inc.Author)
	//log.Fatal("err   ", err, inc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRecord
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
func (m *ItilModel) GetUserByName(name string) (*repository.User, error) {
	stmt := `SELECT id, name, pass  FROM users   
	WHERE name  = ?`
	row := m.db.QueryRow(stmt, name)
	user := &repository.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Pass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRecord
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

			return userID, repository.ErrNoRecord
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
func (m *ItilModel) GetUserRules(id int) (*repository.UserRules, error) {
	stmt := `SELECT id, name  FROM rules   
	WHERE userID  = ?`
	row := m.db.QueryRow(stmt, id)
	userRules := &repository.UserRules{}
	err := row.Scan(&userRules.Id, &userRules.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRecord
		} else {
			return nil, err
		}
	}
	userRules.User = id
	return userRules, nil

}
func (m *ItilModel) GetUserRegionRPC(id int) (string, error) {
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:4041", opts)
	if err != nil {

		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close() // Maybe this should be in a separate function and the error handled?

	c := sample_grpc.NewItilServiceClient(cc)

	// read Region
	fmt.Println("Reading the region")
	readRegionReq := &sample_grpc.GetUserRequest{Id: int32(id)}
	readRegionRes, readRegionErr := c.GetUserRegion(context.Background(), readRegionReq)
	if readRegionErr != nil {
		log.Fatalf("Error happened while reading: %v \n", readRegionErr)
	}

	return readRegionRes.Name, nil
}
