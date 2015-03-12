package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"app/model"
)

func ListComments(w http.ResponseWriter, r *http.Request) {
	comments, err := model.GetComments(
		"WHERE site = ?",
		r.URL.Query().Get("site"),
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Fatal(err)
	}
}
