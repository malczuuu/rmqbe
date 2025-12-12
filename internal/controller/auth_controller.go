package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malczuuu/rmqbe/internal/auth"
)

func NewAuthController(rabbitAuthService auth.RabbitAuthService) AuthController {
	return AuthController{rabbitAuthService: rabbitAuthService}
}

type AuthController struct {
	rabbitAuthService auth.RabbitAuthService
}

func (a *AuthController) User(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	auth := a.rabbitAuthService.User(username, password)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}

func (a *AuthController) Vhost(c *gin.Context) {
	username := c.Query("username")
	vhost := c.Query("vhost")
	ip := c.Query("ip")

	auth := a.rabbitAuthService.Vhost(username, vhost, ip)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}

func (a *AuthController) Resource(c *gin.Context) {
	username := c.Query("username")
	vhost := c.Query("vhost")
	resource := c.Query("resource")
	name := c.Query("name")
	permission := c.Query("permission")

	auth := a.rabbitAuthService.Resource(username, vhost, resource, name, permission)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}

func (a *AuthController) Topic(c *gin.Context) {
	username := c.Query("username")
	vhost := c.Query("vhost")
	resource := c.Query("resource")
	name := c.Query("name")
	permission := c.Query("permission")
	routingKey := c.Query("routing_key")

	auth := a.rabbitAuthService.Topic(username, vhost, resource, name, permission, routingKey)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}
