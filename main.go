package main

import (
	"database/sql"
	"flag"
	"log"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"

	"projectionist/apps"
	"projectionist/config"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			debug.PrintStack()
			log.Fatalln(r)
		}
	}()

	var checker bool
	flag.BoolVar(&checker, "checker", true, "Enable or disable health check")
	flag.Parse()

	log.Printf("Health check mode: %v\n", checker)
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

	syncChan := make(chan string)

	health := apps.NewHealthCkeck(cfg, sqlDB, syncChan)
	if checker {
		err = health.Run()
		if err != nil {
			log.Fatalf("health-check: error: %v", err)
		}
	}

	application, err := apps.NewApp(cfg, sqlDB, syncChan)
	if err != nil {
		log.Fatalf("projectionist-api: error: %v", err)
	}

	application.Run()
}
