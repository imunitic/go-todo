package todos

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("47cc67093475061e3d95369d"))

func List(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	var session *mgo.Session
	var ok bool
	if session, ok = context.Get(req, "session").(*mgo.Session); !ok {
		panic("Data store session not found")
	}

	var result []Todo
	err := session.DB("todos").C("todo").Find(bson.M{"Status": StatusActive}).All(&result)
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
