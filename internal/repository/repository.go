package repository

import (
	"database/sql"
	sqlDB "eduL2_HTTP_BasicAuthServerDB/internal/database"
	mysqlDB "eduL2_HTTP_BasicAuthServerDB/internal/database/mysql"
	"eduL2_HTTP_BasicAuthServerDB/internal/database/sqllitedb"
	"eduL2_HTTP_BasicAuthServerDB/internal/models"
	"log"
)

type Repository struct {
	db  *sql.DB
	drv sqlDB.Driver
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Connect(dbtype, dsn string) error {
	var err error
	log.Println("err  ", dbtype, dsn)
	r.db, err = sql.Open(dbtype, dsn)
	if err != nil {
		return err
	}
	log.Println("err1  ", err)
	if err = r.db.Ping(); err != nil {
		if dbtype == "mysql" {
			//создадим базу Itil

			err = r.db.Close()
			if err != nil {
				log.Println("err112 aasd ", err)
				return err
			} else {
				r.db, err = sql.Open(dbtype, "root:159753@tcp(docker.for.mac.localhost:3306)/")
				if err != nil {

					return err
				}

				err = r.InitBaseRepMain()

				if err != nil {
					return err
				}

				r.db, err = sql.Open(dbtype, dsn)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}
	log.Println("err22  ", err)
	if dbtype == "mysql" {
		r.drv = mysqlDB.NewDriver(r.db)
	} else if dbtype == "sqlite3" {
		r.drv = sqllitedb.NewDriver(r.db)
	} else if dbtype == "postgresql" {

	}
	return err
}
func (r *Repository) StopConnectRep() error {
	var err error
	err = r.db.Close()
	return err
}

func (r *Repository) InsertIncidetRep(name string, author int) (int, error) {
	return r.drv.InsertIncidet(name, author)
}
func (r *Repository) GetIncidentRep(id int) (*models.Incident, error) {
	return r.drv.GetIncident(id)
}
func (r *Repository) InsertUserRep(name, pass string) (int, error) {
	return r.drv.InsertUser(name, pass)
}
func (r *Repository) GetUserByNameRep(name string) (*models.User, error) {
	return r.drv.GetUserByName(name)
}
func (r *Repository) GetUserFirstRep() (int, error) {
	return r.drv.GetUserFirst()
}
func (r *Repository) InsertRulesRep(name string, user int) (int, error) {
	return r.drv.InsertRules(name, user)
}
func (r *Repository) GetUserRulesRep(id int) (*models.UserRules, error) {
	return r.drv.GetUserRules(id)
}
func (r *Repository) InitBaseRep() {
	var usId int
	var err error
	usId = -1
	if usId, err = r.GetUserFirstRep(); usId < 1 {
		if err != nil {
			log.Fatal(err)
			return
		}
		r.drv.InitBase()
		log.Println("Try init db  ", usId)

		//log.Println("Init db  ")
	}
}
func (r *Repository) InitBaseRepMain() error {
	//err := r.drv.InitBaseMain()
	var err error
	_, err = r.db.Exec("CREATE DATABASE IF NOT EXISTS Itil")
	if err != nil {
		return err
	}
	log.Println("CREATE  Itil", err)
	_, err = r.db.Exec("USE Itil")
	if err != nil {

		return err
	}
	log.Println("insert users  ", err)
	_, err = r.db.Exec(`DROP TABLE IF EXISTS incidents;`)
	if err != nil {

		return err
	}

	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS incidents (
	  id int NOT NULL AUTO_INCREMENT,
	  name varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
	  author int NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {

		return err
	}
	log.Println("CREATE  incidents", err)
	_, err = r.db.Exec(`DROP TABLE IF EXISTS rules;`)

	if err != nil {

		return err
	}
	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS rules (
	  id int NOT NULL AUTO_INCREMENT,
	  name varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
	  userID int NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {

		return err
	}
	log.Println("CREATE  rules", err)
	_, err = r.db.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {

		return err
	}
	_, err = r.db.Exec(`CREATE TABLE IF NOT EXISTS users (
	  id int NOT NULL AUTO_INCREMENT,
	  name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
	  pass varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {

		return err
	}
	log.Println("CREATE  users", err)
	return nil

}
