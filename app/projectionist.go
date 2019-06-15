package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"projectionist/provider"

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
	dbProvider  provider.IDBProvider
	sessHandler *session.SessionHandler
}

func NewApp(cfg *config.Config, sqlDB *sql.DB) (*App, error) {
	// initialization database tables
	if err := db.InitTables(sqlDB); err != nil {
		return nil, err
	}

	return &App{
		cfg:         cfg,
		dbProvider:  provider.NewDBProvider(sqlDB),
		sessHandler: session.NewSessionHandler(64, 32),
	}, nil
}

func (a *App) Run() {
	var address = fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	var err error

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
	router.Use(middleware.LoginRequired(a.sessHandler), middleware.AccessControllAllows())

	router.HandleFunc(consts.UrlApiLogin, controllers.LoginApi(a.dbProvider, a.sessHandler)).Methods(http.MethodPost)

	router.HandleFunc(consts.UrlUser, controllers.NewUser(a.dbProvider)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlUser+"/{id}", controllers.GetUser(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUser, controllers.GetUserList(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUser+"/{id}", controllers.UpdateUser(a.dbProvider)).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlUser+"/{id}", controllers.DeleteUser(a.dbProvider)).Methods(http.MethodGet)

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
