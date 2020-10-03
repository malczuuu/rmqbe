package controller

import (
	"fmt"
	"net/http"

	"github.com/malczuuu/rmqbe/internal/auth"
)

func NewAuthController(rabbitAuthService auth.RabbitAuthService) AuthController {
	return AuthController{rabbitAuthService: rabbitAuthService}
}

type AuthController struct {
	rabbitAuthService auth.RabbitAuthService
}

func (c *AuthController) User(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	auth := c.rabbitAuthService.User(username, password)

	result := "deny"
	if auth {
		result = "allow"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func (c *AuthController) Vhost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	vhost := r.FormValue("vhost")
	ip := r.FormValue("ip")

	auth := c.rabbitAuthService.Vhost(username, vhost, ip)

	result := "deny"
	if auth {
		result = "allow"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func (c *AuthController) Resource(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	vhost := r.FormValue("vhost")
	resource := r.FormValue("resource")
	name := r.FormValue("name")
	permission := r.FormValue("permission")

	auth := c.rabbitAuthService.Resource(username, vhost, resource, name, permission)

	result := "deny"
	if auth {
		result = "allow"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

func (c *AuthController) Topic(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	vhost := r.FormValue("vhost")
	resource := r.FormValue("resource")
	name := r.FormValue("name")
	permission := r.FormValue("permission")
	routingKey := r.FormValue("routing_key")

	auth := c.rabbitAuthService.Topic(username, vhost, resource, name, permission, routingKey)

	result := "deny"
	if auth {
		result = "allow"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}
