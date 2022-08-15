package main

import (
	"github.com/Edbeer/Project/config"
	server "github.com/Edbeer/Project/internal/transport/rest"
	"log"

	"github.com/Edbeer/Project/pkg/database/postgres"
)

func main() {
	config := config.GetConfig()

	// postgresql
	psqlClient, err := postgres.NewPsqlDB(config)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v", psqlClient.Stats())
	}
	defer psqlClient.Close()


	s := server.NewServer(config, psqlClient)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
