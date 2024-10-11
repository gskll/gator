package main

import (
	"fmt"
	"log"

	"github.com/gskll/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)

	user := "andrew"
	err = cfg.SetUser(user)
	if err != nil {
		log.Fatal(err)
	}
	cfg2, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg2)
}
