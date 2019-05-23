package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"projectionist/config"
	"projectionist/consts"
	"projectionist/controllers"
	"projectionist/db"
	"projectionist/middleware"
	"projectionist/utils"
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

	log.Printf("<<<<<Projectionst>>>>>")
	log.Printf("Configuration:%+v", cfg)

	var addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	sqlDB, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// initialization database tables
	if err = db.InitTables(sqlDB); err != nil {
		log.Fatal(err)
	}

	// create dir for saving configs files
	if err = utils.CreateDir(consts.PathSaveCfgs); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.Use(middleware.Auth)

	router.HandleFunc(consts.UrlNewUser, controllers.NewUser(sqlDB)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlNewCfg, controllers.NewCfg()).Methods(http.MethodPost)

	log.Printf("Start service on: %v", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
