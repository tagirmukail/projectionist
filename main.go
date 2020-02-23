package main

import (
	"database/sql"
	"flag"
	"github.com/dgraph-io/badger/v2"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	_ "github.com/mattn/go-sqlite3"

	"google.golang.org/grpc/grpclog"

	"projectionist/apps"
	"projectionist/apps/healtchecker"
	"projectionist/config"
)

func init() {
	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)
}

func main() {
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
			grpclog.Fatalln(r)
		}
	}()

	var checker bool
	flag.BoolVar(&checker, "checker", true, "Enable or disable health check")
	flag.Parse()

	grpclog.Infof("Health check mode: %v\n", checker)
	cfg, err := config.NewConfig()
	if err != nil {
		grpclog.Fatalln(err)
	}

	grpclog.Infof("<<<<<Projectionst>>>>>")
	grpclog.Infof("Configuration:%+v", cfg)

	done := make(chan os.Signal, 1)                                    // for graceful down
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // for graceful down

	sqlDB, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer sqlDB.Close()

	badgerDB, err := badger.Open(badger.DefaultOptions("./badgerdb"))
	if err != nil {
		grpclog.Fatalf("open badger db error: %v", err)
	}
	defer badgerDB.Close()

	syncChan := make(chan string, 300)

	health := healtchecker.NewHealthCkeck(cfg, sqlDB, syncChan)
	if checker {
		go func(hc *healtchecker.HealthCheck) {
			err := health.Run()
			if err != nil {
				grpclog.Fatalln(err)
			}
		}(health)
	}

	restApi, err := apps.NewApp(cfg, sqlDB, badgerDB, syncChan)
	if err != nil {
		grpclog.Fatalf("projectionist-api: error: %v", err)
	}

	go restApi.Run()

	go apps.RunGRPC(cfg, sqlDB, badgerDB)

	go apps.RunGrpcApi(cfg)

	<-done // graceful down

	if checker {
		health.Stop()
	}

	grpclog.Infoln("Projectionist is stopped")
}
