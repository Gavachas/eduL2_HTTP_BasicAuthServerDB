package main

import (
	"eduL2_HTTP_BasicAuthServerDB/internal/repository"
	"eduL2_HTTP_BasicAuthServerDB/internal/services"

	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var env string

func main() {
	var dbtype string
	var dsn string

	env = os.Getenv("EnvDBtype")

	if env == "" {
		env = "mysql"
	}
	if env == "mysql" {
		dsn = "root:159753@tcp(docker.for.mac.localhost:3306)/itil?parseTime=true" //для докера
		//dsn = "root:159753@tcp(localhost:3306)/itil?parseTime=true" // локально
		dbtype = "mysql"
	} else if env == "sqlite" {

		dsn = "./internal/database/itild.db"

		//dsn = "./db/itild.db"
		dbtype = "sqlite3"
	} else if env == "postgresql" {

	}
	log.Println("Env  ", env, " dsn", dsn)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	rep := repository.NewRepository()
	err := rep.Connect(dbtype, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		rep.StopConnectRep()
		if err != nil {
			log.Fatal(err)
		}
	}()

	app := services.NewService(errorLog, infoLog, rep)

	srv := &http.Server{
		Addr:     ":4000",
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}
	app.Rep.InitBaseRep()

	log.Println("Starting server on  ", srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)

}
