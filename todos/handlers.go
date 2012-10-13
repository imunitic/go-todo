package todos

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func List(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	defer handlePanic(rw)

	var session *mgo.Session
	var ok bool
	if session, ok = context.Get(req, "session").(*mgo.Session); !ok {
		panic("Data store session not found")
	}

	var result []Todo
	err := session.DB("todos").C("todo").Find(bson.M{"Status": 0}).All(&result)
	if err != nil {
		panic("There are no todos found")
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("Could not serialize result")
	}

	fmt.Fprintf(rw, "%s", b)
}

func Get(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	defer handlePanic(rw)

	vars := mux.Vars(req)
	var session *mgo.Session
	var ok bool
	if session, ok = context.Get(req, "session").(*mgo.Session); !ok {
		panic("Data store session not found")
	}

	result := Todo{}
	err := session.DB("todos").C("todo").FindId(bson.ObjectIdHex(vars["id"])).One(&result)
	if err != nil {
		panic(fmt.Sprintf("Could not find todo with id %s", vars["id"]))
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("Could not serialize result")
	}

	fmt.Fprintf(rw, "%s", b)
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

type jsonError struct {
	Error string `json:"error"`
}

func handlePanic(rw http.ResponseWriter) {
	if r := recover(); r != nil {
		err := fmt.Sprintf("%s", r)
		b, _ := json.Marshal(jsonError{err})
		http.Error(rw, string(b), http.StatusInternalServerError)
	}
}
