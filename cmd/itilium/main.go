package main

import (
	"eduL2_HTTP_BasicAuthServerDB/internal/repository"
	mysqlDB "eduL2_HTTP_BasicAuthServerDB/internal/repository/mysql"
	"eduL2_HTTP_BasicAuthServerDB/internal/repository/pgsql"
	"eduL2_HTTP_BasicAuthServerDB/internal/repository/sqllitedb"
	"eduL2_HTTP_BasicAuthServerDB/internal/services"

	"log"
	"net/http"
	"os"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var env string

func main() {
	var err error
	var dbtype string
	var dsn string
	var drv repository.Driver
	env = os.Getenv("ENVDBTYPE")

	logger, _ := zap.NewProduction()

	if env == "" {
		env = "postgres"
	}
	switch env {
	case "mysql":
		{
			dsn = "root:159753@tcp(docker.for.mac.localhost:3306)/itil?parseTime=true" //для докера
			//dsn = "root:159753@tcp(localhost:3306)/itil?parseTime=true" // локально
			dbtype = "mysql"
			drv, err = mysqlDB.NewDriver(dbtype, dsn)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "sqlite":
		{
			dsn = "./internal/database/itild.db"
			//dsn = "./db/itild.db"
			dbtype = "sqlite3"
			drv, err = sqllitedb.NewDriver(dbtype, dsn)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "postgres":
		{
			//dsn = "postgres:159753@docker.for.mac.localhost/itil?sslmode=disable"
			dsn = "host=docker.for.mac.localhost user=postgres password=159753 dbname=itil sslmode=disable" // для докера
			//dsn = "user=postgres password=159753 dbname=itil sslmode=disable" // локально
			dbtype = "postgres"
			drv, err = pgsql.NewDriver(dbtype, dsn)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	defer func() {
		err = drv.StopConnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Env  ", env, " dsn", dsn)

	//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	rep := repository.NewRepository(drv)

	app := services.NewService(logger, rep)

	srv := &http.Server{
		Addr:     ":4000",
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}
	log.Println("Starting server on  ", srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)

}
