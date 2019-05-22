package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"projectionist/config"
	"projectionist/db"
	"runtime/debug"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
			log.Fatalln(r)
		}
	}()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	var addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	sqlDB, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err = db.CreateTableUsers(sqlDB); err != nil {
		log.Fatal(err)
	}

	if err = db.CreateTableServices(sqlDB); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.Use()

	log.Fatal(http.ListenAndServe(addr, router))
}
