package todos

import (
	"../config"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"labix.org/v2/mgo"
	"net/http"
)

type decorator func(cfg config.Config, f httpHandler) httpHandler
type httpHandler func(rw http.ResponseWriter, req *http.Request)

func H(f func(rw http.ResponseWriter, req *http.Request)) httpHandler {
	return httpHandler(f)
}

func (h httpHandler) D(d decorator, cfg config.Config) httpHandler {
	return d(cfg, h)
}

func Mongo(cfg config.Config, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := mgo.Dial(cfg.Mongo.Url)
		if err != nil {
			panic("Could not initialize data store")
		}

		defer session.Close()
		defer context.Clear(req)

		session.SetMode(mgo.Monotonic, true)

		context.Set(req, "session", session)

		f(rw, req)
	}
}

/*func Authenticate(f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "session")
		if err != nil {
			panic("Unable to create session")
		}

		f(rw, req)
	}
}*/

type jsonError struct {
	Error string `json:"error"`
}

func HandlePanic(cfg config.Config, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Sprintf("%s", r)
				b, _ := json.Marshal(jsonError{err})
				http.Error(rw, string(b), http.StatusInternalServerError)
			}
		}()

		f(rw, req)
	}
}
