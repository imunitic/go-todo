package todos

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"time"
)

var store = sessions.NewCookieStore([]byte("47cc67093475061e3d95369d"))

func List(rw http.ResponseWriter, req *http.Request) {
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user, err := AuthenticatedUser(req)
	if err != nil {
		panic(err)
	}

	var result []Todo
	err = session.DB("todos").C("todo").
		Find(bson.M{"Status": StatusActive, "ow": user.Id.Hex()}).All(&result)
	if err != nil {
		panic(jsonError{"There are no todos found", QueryError})
	}

	Json(rw, result)
}

func Get(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user, err := AuthenticatedUser(req)
	if err != nil {
		panic(err)
	}

	result := Todo{}
	err = session.DB("todos").C("todo").
		Find(bson.M{"_id": bson.ObjectIdHex(vars["id"]), "ow": user.Id.Hex()}).One(&result)
	if err != nil {
		panic(jsonError{fmt.Sprintf("Could not find todo with id %s", vars["id"]), QueryError})
	}

	Json(rw, result)
}

func Delete(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user, err := AuthenticatedUser(req)
	if err != nil {
		panic(err)
	}

	err = session.DB("todos").C("todo").
		Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"]), "ow": user.Id.Hex()})
	if err != nil {
		panic(jsonError{fmt.Sprintf("Could not remove todo with id %s", vars["id"]), QueryError})
	}

	Json(rw, true)
}

func Create(rw http.ResponseWriter, req *http.Request) {
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user, err := AuthenticatedUser(req)
	if err != nil {
		panic(err)
	}

	dueAt, err := time.Parse("2006-01-02 15:04:05", req.FormValue("DueAt"))
	if err != nil {
		dueAt = time.Now().Add(time.Duration(24) * time.Hour)
	}

	priority, err := strconv.Atoi(req.FormValue("Priority"))
	if err != nil {
		priority = 0
	}

	c := session.DB("todos").C("todo")

	id := bson.NewObjectId()
	err = c.Insert(&Todo{Id: id,
		Owner:     user.Id,
		Title:     req.FormValue("Title"),
		Priority:  priority,
		Status:    StatusActive,
		DueAt:     dueAt,
		CreatedAt: time.Now()})
	if err != nil {
		panic(jsonError{"Unable to insert a todo", QueryError})
	}

	Json(rw, struct{ id string }{id: id.Hex()})
}

func Update(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user, err := AuthenticatedUser(req)
	if err != nil {
		panic(err)
	}

	changes := bson.M{}
	hasChanges := false

	if v := req.FormValue("Title"); v != "" {
		changes["ti"] = v
		hasChanges = true
	}

	if v := req.FormValue("DueAt"); v != "" {
		dueAt, err := time.Parse("2006-01-02 15:04:05", v)
		if err == nil {
			changes["du"] = dueAt
			hasChanges = true
		}
	}

	if v := req.FormValue("Priority"); v != "" {
		priority, err := strconv.Atoi(v)
		if err == nil {
			changes["pr"] = priority
			hasChanges = true
		}
	}

	if v := req.FormValue("Status"); v != "" {
		status, err := strconv.Atoi(v)
		if err == nil {
			changes["st"] = status
			hasChanges = true
		}
	}

	if hasChanges {
		err := session.DB("todos").C("todo").
			Update(bson.M{"_id": bson.ObjectIdHex(vars["id"]), "ow": user.Id}, bson.M{"$set": changes})
		if err != nil {
			panic(jsonError{fmt.Sprintf("Could not update todo with id %s", vars["id"]), QueryError})
		}
	}

	Json(rw, hasChanges)
}

func Login(rw http.ResponseWriter, req *http.Request) {
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user := User{}
	err = session.DB("todos").C("user").Find(bson.M{
		"u": req.FormValue("Username"), "p": req.FormValue("Password")}).One(&user)

	if err != nil {
		panic(jsonError{"Authentication failed", AuthenticationError})
	}

	s, err := store.Get(req, "session")
	if err != nil {
		panic(jsonError{"Unable to create session", SessionCreationError})
	}

	s.Values["User"] = user
	s.Save(req, rw)

	Json(rw, true)
}

func Logout(rw http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "session")
	if err != nil {
		panic(jsonError{"Unable to create session", SessionCreationError})
	}

	if _, ok := session.Values["User"]; ok {
		delete(session.Values, "User")
		session.Save(req, rw)
	}

	http.Redirect(rw, req, "/login.html", http.StatusTemporaryRedirect)
}
