package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"labix.org/v2/mgo"
	"net/http"
)

const (
	Unknown int = iota
	MongoSessionCreationError
	SerializationError
	QueryError
	AuthenticationError
	SessionCreationError
)

type jsonError struct {
	ErrorS    string `json:"error"`
	ErrorCode int    `json:"code"`
}

func (j jsonError) Error() string {
	return fmt.Sprintf("%s, %s", j.ErrorS, j.ErrorCode)
}

type decorator func(cfg Config, f httpHandler) httpHandler
type httpHandler func(rw http.ResponseWriter, req *http.Request)

func H(f func(rw http.ResponseWriter, req *http.Request)) httpHandler {
	return httpHandler(f)
}

func (h httpHandler) D(d decorator, cfg Config) httpHandler {
	return d(cfg, h)
}

func Mongo(cfg Config, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := mgo.Dial(cfg.Mongo.Url)
		if err != nil {
			panic(jsonError{"Could not initialize data store", MongoSessionCreationError})
		}

		defer session.Close()
		defer context.Clear(req)

		session.SetMode(mgo.Monotonic, true)

		context.Set(req, "session", session)

		f(rw, req)
	}
}

func Authenticate(cfg Config, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "session")
		if err != nil {
			panic(jsonError{"Unable to create session", SessionCreationError})
		}

		if user, ok := session.Values["User"]; ok {
			defer context.Clear(req)
			context.Set(req, "User", user)
			f(rw, req)
			return
		}

		panic(jsonError{"Permission denied", AuthenticationError})
	}
}

func HandlePanic(cfg Config, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var s = ""
				if err, ok := r.(jsonError); ok {
					b, _ := json.Marshal(err)
					s = string(b)
				} else {
					s = fmt.Sprintf("%s", r)
					b, _ := json.Marshal(jsonError{s, Unknown})
					s = string(b)
				}
				Error(rw, s, http.StatusOK)
			}
		}()

		f(rw, req)
	}
}
