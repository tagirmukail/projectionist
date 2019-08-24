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
	"projectionist/provider"
	"projectionist/utils"
)

type App struct {
	cfg        *config.Config
	dbProvider provider.IDBProvider
}

func NewApp(cfg *config.Config, sqlDB *sql.DB) (*App, error) {
	// initialization database tables
	if err := db.InitTables(sqlDB); err != nil {
		return nil, err
	}

	return &App{
		cfg:        cfg,
		dbProvider: provider.NewDBProvider(sqlDB),
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
	router.Use(
		middleware.JwtAuthentication(a.cfg.TokenSecretKey),
		middleware.AccessControllAllows(a.cfg.AccessAddresses),
	)

	router.HandleFunc(consts.UrlApiLoginV1, controllers.LoginApi(a.dbProvider, a.cfg.TokenSecretKey)).Methods(http.MethodPost)

	router.HandleFunc(consts.UrlUserV1, controllers.NewUser(a.dbProvider)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlUserV1+"/{id}", controllers.GetUser(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUserV1, controllers.GetUserList(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUserV1+"/{id}", controllers.UpdateUser(a.dbProvider)).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlUserV1+"/{id}", controllers.DeleteUser(a.dbProvider)).Methods(http.MethodGet)

	router.HandleFunc(consts.UrlCfgV1, controllers.NewCfg()).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlCfgV1+"/{id}", controllers.GetCfg()).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlCfgV1, controllers.GetCfgList()).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlCfgV1+"/{id}", controllers.UpdateCfg()).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlCfgV1+"/{id}", controllers.DeleteCfg()).Methods(http.MethodDelete)

	router.HandleFunc(consts.UrlServiceV1, controllers.NewService()).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlServiceV1+"/{id}", controllers.GetService()).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlServiceV1, controllers.GetServiceList()).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlServiceV1+"/{id}", controllers.UpdateService()).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlServiceV1+"/{id}", controllers.DeleteService()).Methods(http.MethodDelete)

	return router
}
