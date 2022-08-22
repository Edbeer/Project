package main

import (
	"github.com/Edbeer/Project/config"
	server "github.com/Edbeer/Project/internal/transport/rest"

	"github.com/Edbeer/Project/pkg/database/postgres"
	"github.com/Edbeer/Project/pkg/database/redis"
	"github.com/Edbeer/Project/pkg/logger"
)

func main() {
	// init config
	config := config.GetConfig()

	// init logger
	logger := logger.NewApiLogger(config)
	logger.InitLogger()

	// postgresql
	psqlClient, err := postgres.NewPsqlDB(config)
	if err != nil {
		logger.Fatalf("Postgresql init: %s", err)
	} else {
		logger.Infof("Postgres connected, Status: %#v", psqlClient.Stats())
	}
	defer psqlClient.Close()

	// redis
	redisClient := redis.NewRedisClient(config)
	defer redisClient.Close()
	logger.Info("Redis connetcted")
	
	logger.Info("Starting auth server")
	s := server.NewServer(config, psqlClient, redisClient, logger)
	if err := s.Run(); err != nil {
		logger.Fatal(err)
	}
}
