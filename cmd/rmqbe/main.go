package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/malczuuu/rmqbe/internal/apidoc"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})

	addr := "0.0.0.0:8000"

	log.WithField("addr", addr).Info("HTTP server is being initialized")
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")

	router.HandleFunc("/user", userHandler).Methods("POST")
	router.HandleFunc("/vhost", vhostHandler).Methods("POST")
	router.HandleFunc("/resource", resourceHandler).Methods("POST")
	router.HandleFunc("/topic", topicHandler).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	docs := apidoc.GetStructure()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	result := "deny"

	if username != "" && password != "" {
		result = "allow"
	}

	log.WithFields(
		map[string]interface{}{
			"username": username,
			"result":   result,
		}).Info("Authenticate user by username and password")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func vhostHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	vhost := r.FormValue("vhost")
	ip := r.FormValue("ip")

	result := "deny"

	if vhost == "/" && username != "" {
		result = "allow"
	}

	log.WithFields(
		map[string]interface{}{
			"username": username,
			"vhost":    vhost,
			"ip":       ip,
			"result":   result,
		}).Info("Authorize user to virtual host")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	vhost := r.FormValue("vhost")
	resource := r.FormValue("resource")
	name := r.FormValue("name")
	permission := r.FormValue("permission")

	result := "allow"

	log.WithFields(
		map[string]interface{}{
			"username":   username,
			"vhost":      vhost,
			"resource":   resource,
			"name":       name,
			"permission": permission,
			"result":     result,
		}).Info("Authorize user to resource")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	vhost := r.FormValue("vhost")
	resource := r.FormValue("resource")
	name := r.FormValue("name")
	permission := r.FormValue("permission")
	routingKey := r.FormValue("routing_key")

	result := "allow"

	log.WithFields(
		map[string]interface{}{
			"username":    username,
			"vhost":       vhost,
			"resource":    resource,
			"name":        name,
			"permission":  permission,
			"routing_key": routingKey,
			"result":      result,
		}).Info("Authorize user to topic")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}
