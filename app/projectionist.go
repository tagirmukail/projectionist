package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"projectionist/config"
	"projectionist/consts"
	"projectionist/controllers"
	"projectionist/db"
	"projectionist/middleware"
	"projectionist/session"
	"projectionist/utils"
)

type App struct {
	cfg         *config.Config
	db          *sql.DB
	sessHandler *session.SessionHandler
}

func NewApp(cfg *config.Config, db *sql.DB) *App {
	return &App{
		cfg:         cfg,
		db:          db,
		sessHandler: session.NewSessionHandler(64, 32),
	}
}

func (a *App) Run() {
	var address = fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	var err error

	// initialization database tables
	if err = db.InitTables(a.db); err != nil {
		log.Fatal(err)
	}

	// create dir for saving configs files
	if err = utils.CreateDir(consts.PathSaveCfgs); err != nil {
		log.Fatal(err)
	}

	router := a.newRouter()

	log.Printf("Start service on: %v", address)
	log.Fatal(http.ListenAndServe(address, router))
}

func (a *App) newRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LoginRequired(a.sessHandler))

	router.HandleFunc(consts.UrlApiLogin, controllers.LoginApi(a.db, a.sessHandler)).Methods(http.MethodPost)

	router.HandleFunc(consts.UrlUser, controllers.NewUser(a.db)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlUser+"/{id}", controllers.GetUser(a.db)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUser, controllers.GetUserList(a.db)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUser+"/{id}", controllers.UpdateUser(a.db)).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlUser+"/{id}", controllers.DeleteUser(a.db)).Methods(http.MethodGet)

	router.HandleFunc(consts.UrlCfg, controllers.NewCfg()).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlCfg+"/{id}", controllers.GetCfg()).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlCfg, controllers.GetCfgList()).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlCfg+"/{id}", controllers.UpdateCfg()).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlCfg+"/{id}", controllers.DeleteCfg()).Methods(http.MethodDelete)

	router.HandleFunc(consts.UrlService, controllers.NewService()).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlService+"/{id}", controllers.GetService()).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlService, controllers.GetServiceList()).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlService+"/{id}", controllers.UpdateService()).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlService+"/{id}", controllers.DeleteService()).Methods(http.MethodDelete)

	return router
}
