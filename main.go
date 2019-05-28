package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"projectionist/config"
	"projectionist/consts"
	"projectionist/controllers"
	"projectionist/db"
	"projectionist/middleware"
	"projectionist/models"
	"projectionist/session"
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

	user := models.User{}
	usersNotEmpty, err := user.TableNotEmpty(sqlDB)
	if err != nil {
		log.Fatalln(err)
	}

	sessHandl := session.NewSessionHandler(64, 32)

	t := make(map[string]*template.Template)
	t["login.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))
	t["services-index.html"] = template.Must(
		template.ParseFiles(
			"templates/base.html",
			"templates/services-index.html",
		),
	)

	router := mux.NewRouter()
	router.Use(middleware.FirstAuth(&usersNotEmpty), middleware.LoginRequired(sessHandl))

	router.HandleFunc(consts.UrlApiLogin, controllers.LoginApi(sqlDB)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlNewUser, controllers.NewUser(sqlDB, &usersNotEmpty)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlNewCfg, controllers.NewCfg()).Methods(http.MethodPost)

	router.HandleFunc("/login", controllers.Login(t["login.html"])).Methods(http.MethodGet)
	router.HandleFunc("/", controllers.ServicesIndex(t["services-index.html"]))

	log.Printf("Start service on: %v", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
