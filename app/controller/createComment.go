package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"app/model"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	comment := &model.Comment{}
	if err := json.NewDecoder(r.Body).Decode(comment); err != nil {
		http.Error(w, "unable to parse the comment", 400)
		return
	}

	if err := model.CreateComment(comment); err != nil {
		http.Error(w, "unable to save the comment "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		log.Fatal(err)
	}
}
