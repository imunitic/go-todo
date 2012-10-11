package main

import (
	"./todos"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/todos", todos.List).Methods("GET")
	r.HandleFunc("/todo", todos.Create).Methods("PUT")
	r.HandleFunc("/todo/{id}", todos.Get).Methods("GET")
	r.HandleFunc("/todo/{id}", todos.Update).Methods("POST")
	r.HandleFunc("/todo/{id}", todos.Delete).Methods("DELETE")

	http.Handle("/api/", r)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./html"))))
	http.ListenAndServe(":8080", nil)
}
