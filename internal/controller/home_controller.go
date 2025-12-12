package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malczuuu/rmqbe/internal/apidoc"
)

func NewHomeController() HomeController {
	return HomeController{}
}

type HomeController struct{}

func (h *HomeController) Home(c *gin.Context) {
	docs := apidoc.GetStructure()
	c.JSON(http.StatusOK, docs)
}
