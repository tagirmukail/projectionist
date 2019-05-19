package main

import (
	"fmt"
	"log"
	"projectionist/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v", cfg)
}
