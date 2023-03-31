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

var env_dbtype, env_dsn string

func main() {
	var err error

	var drv repository.Driver
	env_dbtype = os.Getenv("ENVDBTYPE")
	env_dsn = os.Getenv("ENVDSN")

	if env_dbtype == "" {
		env_dbtype = "postgres"
	}

	switch env_dbtype {
	case "mysql":
		{
			if env_dsn == "" {
				env_dsn = "root:159753@tcp(docker.for.mac.localhost:3306)/itil?parseTime=true" // для докера
				//env_dsn = "root:159753@tcp(localhost:3306)/itil?parseTime=true" // локально
			}
			drv, err = mysqlDB.NewDriver(env_dbtype, env_dsn)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "sqlite3":
		{
			if env_dsn == "" {
				env_dsn = "./internal/database/itild.db" // для докера
				//env_dsn = "./db/itild.db" // локально
			}
			drv, err = sqllitedb.NewDriver(env_dbtype, env_dsn)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "postgres":
		{
			if env_dsn == "" {
				env_dsn = "host=docker.for.mac.localhost user=postgres password=159753 dbname=itil sslmode=disable" // для докера
				//env_dsn = "user=postgres password=159753 dbname=itil sslmode=disable" // локально
			}
			drv, err = pgsql.NewDriver(env_dbtype, env_dsn)
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
	log.Println("Env  ", env_dbtype, " dsn", env_dsn)

	logger, _ := zap.NewProduction()
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
