package controller

import (
	"encoding/json"
	"net/http"

	"github.com/malczuuu/rmqbe/internal/apidoc"
)

func NewHomeController() HomeController {
	return HomeController{}
}

type HomeController struct{}

func (c *HomeController) Home(w http.ResponseWriter, r *http.Request) {
	docs := apidoc.GetStructure()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
