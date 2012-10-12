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

	r.HandleFunc("/todos", todos.List).Methods("GET")
	r.HandleFunc("/todo", todos.Create).Methods("PUT")
	r.HandleFunc("/todo/{id}", todos.Get).Methods("GET")
	r.HandleFunc("/todo/{id}", todos.Update).Methods("POST")
	r.HandleFunc("/todo/{id}", todos.Delete).Methods("DELETE")

	http.Handle("/api/", r)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./html"))))
	http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.Port), nil)
}
