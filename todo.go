package main

import (
	"./config"
	"./todos"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	cfg := config.Config{}
	if err := cfg.Load("./config.json"); err != nil {
		panic(err)
	}

	logout := todos.H(todos.Logout).D(todos.HandlePanic, cfg)
	login := todos.H(todos.Login).D(todos.Mongo, cfg).D(todos.HandlePanic, cfg)
	list := todos.H(todos.List).D(todos.Mongo, cfg).D(todos.Authenticate, cfg).D(todos.HandlePanic, cfg)
	create := todos.H(todos.Create).D(todos.Mongo, cfg).D(todos.Authenticate, cfg).D(todos.HandlePanic, cfg)
	get := todos.H(todos.Get).D(todos.Mongo, cfg).D(todos.Authenticate, cfg).D(todos.HandlePanic, cfg)
	update := todos.H(todos.Update).D(todos.Mongo, cfg).D(todos.Authenticate, cfg).D(todos.HandlePanic, cfg)
	delete := todos.H(todos.Delete).D(todos.Mongo, cfg).D(todos.Authenticate, cfg).D(todos.HandlePanic, cfg)

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
