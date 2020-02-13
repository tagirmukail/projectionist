package apps

import (
	"database/sql"
	"fmt"
	"github.com/dgraph-io/badger"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"google.golang.org/grpc/grpclog"

	"projectionist/config"
	"projectionist/consts"
	"projectionist/controllers"
	"projectionist/db"
	"projectionist/middleware"
	"projectionist/provider"
)

type App struct {
	syncChan    chan string
	cfg         *config.Config
	dbProvider  provider.IDBProvider
	cfgProvider provider.IDBProvider
}

func NewApp(cfg *config.Config, sqlDB *sql.DB, badgerDB *badger.DB, syncShan chan string) (*App, error) {
	// initialization database tables
	if err := db.InitTables(sqlDB); err != nil {
		return nil, err
	}

	cfgProvider, err := provider.NewCfgProvider(badgerDB)
	if err != nil {
		return nil, err
	}

	return &App{
		syncChan:    syncShan,
		cfg:         cfg,
		dbProvider:  provider.NewDBProvider(sqlDB),
		cfgProvider: cfgProvider,
	}, nil
}

func (a *App) Run() {
	var address = fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)

	router := a.newRouter()

	grpclog.Infof("Start rest api service on: %v", address)
	grpclog.Fatal(http.ListenAndServe(
		address,
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins(a.cfg.AccessAddresses),
		)(router)),
	)
}

func (a *App) newRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(
		middleware.JwtAuthentication(a.cfg.TokenSecretKey),
	)

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./apps/swagger/")))
	router.PathPrefix("/swaggerui/").Handler(sh)

	router.HandleFunc(consts.UrlApiLoginV1, controllers.LoginApi(a.dbProvider, a.cfg.TokenSecretKey)).Methods(http.MethodPost)

	router.HandleFunc(consts.UrlUserV1, controllers.NewUser(a.dbProvider)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlUserV1, controllers.GetUserList(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUserV1+"/{id}", controllers.GetUser(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlUserV1+"/{id}", controllers.UpdateUser(a.dbProvider)).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlUserV1+"/{id}", controllers.DeleteUser(a.dbProvider)).Methods(http.MethodGet)

	router.HandleFunc(consts.UrlCfgV1, controllers.NewCfg(a.cfgProvider)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlCfgV1, controllers.GetCfgList(a.cfgProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlCfgV1+"/{id}", controllers.GetCfg(a.cfgProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlCfgV1+"/{id}", controllers.UpdateCfg(a.cfgProvider)).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlCfgV1+"/{id}", controllers.DeleteCfg(a.cfgProvider)).Methods(http.MethodDelete)

	router.HandleFunc(consts.UrlServiceV1, controllers.NewService(a.dbProvider, a.syncChan)).Methods(http.MethodPost)
	router.HandleFunc(consts.UrlServiceV1, controllers.GetServiceList(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlServiceV1+"/{id}", controllers.GetService(a.dbProvider)).Methods(http.MethodGet)
	router.HandleFunc(consts.UrlServiceV1+"/{id}", controllers.UpdateService(a.dbProvider, a.syncChan)).Methods(http.MethodPut)
	router.HandleFunc(consts.UrlServiceV1+"/{id}", controllers.DeleteService(a.dbProvider, a.syncChan)).Methods(http.MethodDelete)

	return router
}
