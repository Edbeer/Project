package main

import (
	"log"

	"github.com/Edbeer/Project/config"
	server "github.com/Edbeer/Project/internal/transport/rest"

	"github.com/Edbeer/Project/pkg/database/postgres"
	"github.com/Edbeer/Project/pkg/database/redis"
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

	// redis
	redisClient := redis.NewRedisClient(config)
	defer redisClient.Close()
	log.Println("Redis connetcted")
	
	s := server.NewServer(config, psqlClient, redisClient)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
