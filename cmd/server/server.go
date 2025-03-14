package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gophKeeper/internal/server/config"
	v2 "gophKeeper/internal/server/controller/http/v1"
	"gophKeeper/internal/server/db"
	"gophKeeper/internal/server/middleware"
	repository1 "gophKeeper/internal/server/services/auth/repository"
	usecase1 "gophKeeper/internal/server/services/auth/usecase"
	repository3 "gophKeeper/internal/server/services/lockbox/repository"
	usecase3 "gophKeeper/internal/server/services/lockbox/usecase"
	repository2 "gophKeeper/internal/server/services/users/repository"
	usecase2 "gophKeeper/internal/server/services/users/usecase"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	logger, _ := zap.NewProduction()
	gin.SetMode(cfg.App.Mode)

	database, err := db.Init(ctx, cfg)
	if err != nil {
		logger.Error("Error database: " + err.Error())
		return
	}

	if cfg.App.Mode == "debug" {
		if err := database.CheckMigrations(ctx, cfg.Pg.Migrate); err != nil {
			logger.Error("Error database migration error: " + err.Error())
			return
		}
	}

	mware := middleware.NewMiddlewareService(cfg, database)

	authRepos := repository1.NewAuthRepository(database)
	authUsecase := usecase1.NewAuthUsecase(authRepos)

	userRepos := repository2.NewUserRepository(database)
	userUsecase := usecase2.NewUserUsecase(userRepos)

	lockBoxRepos := repository3.NewLockBoxRepo(database)
	lockBoxUsecase := usecase3.NewLockBoxUsecase(lockBoxRepos)

	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				lockBoxRepos.PurgeExpiredLocks(ctx)
			case <-ctx.Done():
				return
			}

		}
	}()

	router.Use(cors.New(corsConfig))
	api := router.Group("/api")
	{
		v2.NewAuthHandler(cfg, api, authUsecase, mware)
		v2.NewLockBoxHandlerHandler(cfg, api, lockBoxUsecase, mware)
		v2.NewUserHandler(cfg, api, userUsecase, mware)
	}

	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
		Addr:         cfg.App.Host + ":" + cfg.App.Port,
		Handler:      router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}

}
