package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	cfg := Config{}
	if err := cfg.Load("./config.json"); err != nil {
		panic(err)
	}

	logout := H(Logout).D(HandlePanic, cfg)
	login := H(Login).D(Mongo, cfg).D(HandlePanic, cfg)
	list := H(List).D(Mongo, cfg).D(Authenticate, cfg).D(HandlePanic, cfg)
	create := H(Create).D(Mongo, cfg).D(Authenticate, cfg).D(HandlePanic, cfg)
	get := H(Get).D(Mongo, cfg).D(Authenticate, cfg).D(HandlePanic, cfg)
	update := H(Update).D(Mongo, cfg).D(Authenticate, cfg).D(HandlePanic, cfg)
	delete := H(Delete).D(Mongo, cfg).D(Authenticate, cfg).D(HandlePanic, cfg)

	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")
	r.HandleFunc("/todos", list).Methods("GET")
	r.HandleFunc("/todo", create).Methods("PUT")
	r.HandleFunc("/todo/{id}", get).Methods("GET")
	r.HandleFunc("/todo/{id}", update).Methods("POST")
	r.HandleFunc("/todo/{id}", delete).Methods("DELETE")

	http.Handle("/api/", r)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(cfg.Server.WebRoot))))
	http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.Port), nil)
}
