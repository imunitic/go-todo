package todos

import (
	"../config"
	"github.com/gorilla/context"
	"labix.org/v2/mgo"
	"net/http"
)

type httpHandler func(rw http.ResponseWriter, req *http.Request)

func DataStore(cfg config.Mongo, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := mgo.Dial(cfg.Url)
		if err != nil {
			http.Error(rw, "Could not initialize data store", http.StatusInternalServerError)
			return
		}

		defer session.Close()
		defer context.Clear(req)

		session.SetMode(mgo.Monotonic, true)

		context.Set(req, "session", session)

		f(rw, req)
	}
}
