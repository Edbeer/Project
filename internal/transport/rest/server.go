package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/service"
	"github.com/Edbeer/Project/internal/storage/psql"
	"github.com/Edbeer/Project/internal/storage/redis"
	"github.com/Edbeer/Project/internal/transport/rest/api"
	"github.com/Edbeer/Project/pkg/hash"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	psql   *sqlx.DB
	redis  *redis.Client
}

// New Server constructor
func NewServer(config *config.Config, psql *sqlx.DB, redis *redis.Client) *Server {
	return &Server{
		echo:   echo.New(),
		config: config,
		psql:   psql,
		redis:  redis,
	}
}

func (s *Server) Run() error {
	// Services, Repos & API Handlers
	config := config.GetConfig()

	hash := hash.NewSHA1Hasher()

	psql := psql.NewStorage(s.psql)
	redis := redisrepo.NewStorage(s.redis)
	service := service.NewServices(service.Deps{
		Config:       s.config,
		PsqlStorage:  psql,
		RedisStorage: redis,
		Hash:         hash,
	})
	handlers := api.NewHandlers(api.Deps{
		UserService: service.User,
		Config:      config,
	})
	if err := handlers.Init(s.echo); err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:           s.config.Server.Port,
		ReadTimeout:    time.Second * time.Duration(s.config.Server.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.config.Server.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("Error starting Server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.echo.Server.Shutdown(ctx)
}
