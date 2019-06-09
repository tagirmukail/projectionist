package main

import (
	"database/sql"
	"log"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"

	"projectionist/app"
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

	application, err := app.NewApp(cfg, sqlDB)
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
