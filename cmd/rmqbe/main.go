package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/malczuuu/rmqbe/internal/auth"
	"github.com/malczuuu/rmqbe/internal/config"
	"github.com/malczuuu/rmqbe/internal/controller"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetLevel(log.DebugLevel)

	config := config.ReadConfig()

	log.WithField("config", config).Info("Starting RMQ BE")

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	client.Database(config.MongoDatabase).CreateCollection(ctx, config.MongoUsersCollection)

	auth := auth.NewRabbitAuthService(client, config)

	homeController := controller.NewHomeController()
	authController := controller.NewAuthController(auth)

	addr := "0.0.0.0:8000"

	log.WithField("addr", addr).Info("HTTP server is being initialized")
	router := mux.NewRouter()
	router.HandleFunc("/", homeController.Home).Methods("GET")

	router.HandleFunc("/user", authController.User).Methods("POST")
	router.HandleFunc("/vhost", authController.Vhost).Methods("POST")
	router.HandleFunc("/resource", authController.Resource).Methods("POST")
	router.HandleFunc("/topic", authController.Topic).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
