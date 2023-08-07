package app

import (
	"github.com/atlant1da-404/droplet/config"
	controller "github.com/atlant1da-404/droplet/internal/controller/http"
	"github.com/atlant1da-404/droplet/internal/entity"
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/internal/storage"
	"github.com/atlant1da-404/droplet/pkg/auth"
	"github.com/atlant1da-404/droplet/pkg/database"
	"github.com/atlant1da-404/droplet/pkg/hash"
	"github.com/atlant1da-404/droplet/pkg/httpserver"
	"github.com/atlant1da-404/droplet/pkg/logger"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	sql, err := database.NewPostgreSQL(database.PostgreSQLConfig{
		User:     cfg.PostgreSQL.User,
		Password: cfg.PostgreSQL.Password,
		Host:     cfg.PostgreSQL.Host,
		Database: cfg.PostgreSQL.Database,
	})
	if err != nil {
		log.Fatal("failed to init postgresql", "err", err)
	}

	err = sql.DB.AutoMigrate(
		&entity.User{},
		&entity.Account{},
		&entity.AccountDevices{},
		&entity.AccountSettings{},
	)
	if err != nil {
		log.Fatal("automigration failed", "err", err)
	}

	storages := service.Storages{
		UserStorage:    storage.NewUserStorage(sql),
		AccountStorage: storage.NewAccountStorage(sql),
	}

	databases := map[string]database.Database{
		"postgreSQL": sql,
	}

	serviceOptions := &service.Options{
		Storages: &storages,
		Config:   cfg,
		Logger:   log,
		Hash:     hash.NewHash(),
		Auth:     auth.NewAuth(),
	}

	services := service.Services{
		AuthService:    service.NewAuthService(serviceOptions),
		AccountService: service.NewAccountService(serviceOptions),
	}

	httpHandler := gin.New()

	controller.New(&controller.Options{
		Handler:  httpHandler,
		Services: services,
		Logger:   log,
		Config:   cfg,
	})

	httpServer := httpserver.New(
		httpHandler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Second*60),
		httpserver.WriteTimeout(time.Second*60),
		httpserver.ShutdownTimeout(time.Second*30),
	)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify", "err", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown", "err", err)
	}

	for _, db := range databases {
		err = db.Close()
		if err != nil {
			log.Error("app - Run - db.Close", "err", err)
		}
	}
}
