package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"app/controller"
	"app/model"
	"app/util"
)

func main() {
	model.Init()
	defer model.Deinit()

	r := mux.NewRouter()

	r.HandleFunc("/comments", controller.ListComments).
		Methods("GET").
		Queries("site", "{site}")

	r.HandleFunc("/comment", controller.CreateComment).
		Methods("POST")

	r.HandleFunc("/comment", controller.EditComment).
		Methods("PUT")

	r.HandleFunc("/comment/{id:[0-9+]}", controller.DeleteComment).
		Methods("DELETE")

	http.ListenAndServe("127.0.0.1:8080", util.Middleware(r))
}
