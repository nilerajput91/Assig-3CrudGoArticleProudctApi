package controllers

import (
	"net/http"

	"github.com/nilerajput91/Assig-3CrudGoArticleProudctApi/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
