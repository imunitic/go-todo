package todos

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func List(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "Fetching todos")
}

func Get(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Fprintf(rw, "Fetching todo with id %s", vars["id"])
}

func Delete(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Fprintf(rw, "Deleteing todo with id %s", vars["id"])

}

func Create(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "Creating a new todo")

}

func Update(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Fprintf(rw, "Updating todo with id %s", vars["id"])
}
