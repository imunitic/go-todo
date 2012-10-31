package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"labix.org/v2/mgo"
	"net/http"
)

func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

func Json(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, err := json.Marshal(data)
	if err != nil {
		panic(jsonError{"Could not serialize result", SerializationError})
	}

	fmt.Fprintln(w, fmt.Sprintf("%s", b))
}

func MongoSession(req *http.Request) (s *mgo.Session, err error) {
	var ok bool
	if s, ok = context.Get(req, "session").(*mgo.Session); !ok {
		err = jsonError{"Data store session not found", MongoSessionCreationError}
	}

	return
}

func AuthenticatedUser(req *http.Request) (user User, err error) {
	var ok bool
	user, ok = context.Get(req, "User").(User)
	if !ok {
		err = jsonError{"Invalid user data", AuthenticationError}
	}

	return
}
