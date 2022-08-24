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
	"github.com/Edbeer/Project/pkg/jwt"
	"github.com/Edbeer/Project/pkg/logger"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	psql   *sqlx.DB
	redis  *redis.Client
	logger  logger.Logger
}

// New Server constructor
func NewServer(config *config.Config, psql *sqlx.DB, redis *redis.Client, logger logger.Logger) *Server {
	return &Server{
		echo:   echo.New(),
		config: config,
		psql:   psql,
		redis:  redis,
		logger: logger,
	}
}

func (s *Server) Run() error {
	// Services, Repos & API Handlers
	hash := hash.NewSHA1Hasher()

	tokenManager, err := jwt.NewManager(s.config.Server.JwtSecretKey)
	if err != nil {
		return err
	}
	psql := psql.NewStorage(s.psql)
	redis := redisrepo.NewStorage(redisrepo.Deps{
		Redis: s.redis,
		Manager: tokenManager,
	})
	service := service.NewServices(service.Deps{
		Config:       s.config,
		PsqlStorage:  psql,
		RedisStorage: redis,
		Hash:         hash,
		TokenManager: tokenManager,
	})
	handlers := api.NewHandlers(api.Deps{
		UserService:    service.User,
		SessionService: service.Session,
		Config:         s.config,
	})
	if err := handlers.Init(s.echo, s.logger); err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:           s.config.Server.Port,
		ReadTimeout:    time.Second * time.Duration(s.config.Server.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.config.Server.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		s.logger.Infof("Server is listening on port: %s", s.config.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: %v", err)

		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
