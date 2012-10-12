package todos

import (
	"github.com/gorilla/context"
	"labix.org/v2/mgo"
	"net/http"
)

type httpHandler func(rw http.ResponseWriter, req *http.Request)

func DecorateWithMongoSession(url string, f httpHandler) httpHandler {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := mgo.Dial(url)
		if err != nil {
			http.Error(rw, "Mongo session", http.StatusInternalServerError)
			return
		}

		defer session.Close()
		defer context.Clear(req)

		context.Set(req, "session", session)

		f(rw, req)
	}
}
