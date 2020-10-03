package main

import (
	"context"
	"net/http"
	"time"

	"github.com/malczuuu/rmqbe/internal/controller"

	"github.com/malczuuu/rmqbe/internal/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	client.Database("rmqbe").CreateCollection(ctx, "rmqbe")

	authManager := auth.NewAuthManager(client)

	homeController := controller.NewHomeController()
	authController := controller.NewAuthController(authManager)

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
