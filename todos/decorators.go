package todos

import (
	"github.com/gorilla/context"
	"labix.org/v2/mgo"
	"net/http"
)

func DecorateWithMongoSession(url string, f func(rw http.ResponseWriter, req *http.Request)) func(rw http.ResponseWriter, req *http.Request) {
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
