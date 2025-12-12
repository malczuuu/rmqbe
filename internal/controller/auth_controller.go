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
	username := c.PostForm("username")
	password := c.PostForm("password")

	auth := a.rabbitAuthService.User(username, password)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}

func (a *AuthController) Vhost(c *gin.Context) {
	username := c.PostForm("username")
	vhost := c.PostForm("vhost")
	ip := c.PostForm("ip")

	auth := a.rabbitAuthService.Vhost(username, vhost, ip)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}

func (a *AuthController) Resource(c *gin.Context) {
	username := c.PostForm("username")
	vhost := c.PostForm("vhost")
	resource := c.PostForm("resource")
	name := c.PostForm("name")
	permission := c.PostForm("permission")

	auth := a.rabbitAuthService.Resource(username, vhost, resource, name, permission)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}

func (a *AuthController) Topic(c *gin.Context) {
	username := c.PostForm("username")
	vhost := c.PostForm("vhost")
	resource := c.PostForm("resource")
	name := c.PostForm("name")
	permission := c.PostForm("permission")
	routingKey := c.PostForm("routing_key")

	auth := a.rabbitAuthService.Topic(username, vhost, resource, name, permission, routingKey)

	result := "deny"
	if auth {
		result = "allow"
	}

	c.String(http.StatusOK, result)
}
