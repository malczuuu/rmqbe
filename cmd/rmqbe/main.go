package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/malczuuu/rmqbe/internal/auth"
	"github.com/malczuuu/rmqbe/internal/config"
	"github.com/malczuuu/rmqbe/internal/controller"
	"github.com/malczuuu/rmqbe/internal/logging"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	config := config.ReadConfig()
	logging.ConfigureLogger(&config)

	log.Info().Msg("starting rmqbe service")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		panic(err)
	}

	client.Database(config.MongoDatabase).CreateCollection(ctx, config.MongoUsersCollection)

	auth := auth.NewRabbitAuthService(client, config)

	homeController := controller.NewHomeController()
	authController := controller.NewAuthController(auth)

	addr := "0.0.0.0:8000"

	log.Info().Str("addr", addr).Msg("HTTP server is being initialized")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", homeController.Home)

	router.POST("/user", authController.User)
	router.POST("/vhost", authController.Vhost)
	router.POST("/resource", authController.Resource)
	router.POST("/topic", authController.Topic)

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Info().Str("addr", addr).Msg("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server exited with error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info().Str("signal", sig.String()).Msg("commencing graceful shutdown")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("graceful shutdown completed")
}
