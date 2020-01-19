package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"runtime/debug"
	"sync"

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
			log.Fatalln(r)
		}
	}()

	var grpc bool
	var checker bool
	flag.BoolVar(&checker, "checker", true, "Enable or disable health check")
	flag.BoolVar(&grpc, "grpc", true, "Enable or disable grpc")
	flag.Parse()

	log.Printf("Health check mode: %v\n", checker)
	log.Printf("GRPC: %v\n", grpc)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("<<<<<Projectionst>>>>>")
	log.Printf("Configuration:%+v", cfg)

	sqlDB, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	syncChan := make(chan string, 300)

	health := healtchecker.NewHealthCkeck(cfg, sqlDB, syncChan)
	if checker {
		go health.Run()
	}

	wg := &sync.WaitGroup{}

	application, err := apps.NewApp(cfg, sqlDB, syncChan)
	if err != nil {
		log.Fatalf("projectionist-api: error: %v", err)
	}

	wg.Add(1)
	go apps.RunGRPC(wg, cfg, sqlDB)

	wg.Add(1)
	go application.Run(wg)

	wg.Add(1)
	go apps.RunGrpcApi(wg, cfg)

	wg.Wait()

	health.Stop()
}
