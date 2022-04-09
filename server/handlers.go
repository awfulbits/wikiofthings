package server

import (
	"net/http"

	"github.com/awfulbits/wikiofthings/database"
	"github.com/gorilla/mux"
)

func titleHandler(db *database.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		titlePage, err := loadTitle(vars["title"], db)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		renderTitleTemplate(w, "title", titlePage)
	})
}
