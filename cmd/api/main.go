package main

import (
	"github.com/Edbeer/Project/config"
	server "github.com/Edbeer/Project/internal/transport/rest"
	"github.com/opentracing/opentracing-go"
	jLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/Edbeer/Project/pkg/database/postgres"
	"github.com/Edbeer/Project/pkg/database/redis"
	"github.com/Edbeer/Project/pkg/logger"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
)

// @title           Auth App Api
// @version         1.0
// @description     This is an example of Auth

// @BasePath /api/

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
	
	jaegerCfgInstance := jConfig.Configuration{
		ServiceName: "REST_API",
		Sampler: &jConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jConfig.ReporterConfig{
			LogSpans:           config.Jaeger.LogSpans,
			LocalAgentHostPort: config.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jConfig.Logger(jLog.StdLogger),
		jConfig.Metrics(metrics.NullFactory),
	)
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	logger.Info("Opentracing connected")

	logger.Info("Starting auth server")
	s := server.NewServer(config, psqlClient, redisClient, logger)
	if err := s.Run(); err != nil {
		logger.Fatal(err)
	}
}
