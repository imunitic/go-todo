package todos

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo/bson"
	"net/http"
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
	err = session.DB("todos").C("todo").Find(bson.M{"Status": StatusActive, "Owner": user.Id.Hex()}).All(&result)
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
	err = session.DB("todos").C("todo").Find(bson.M{"_id": bson.ObjectIdHex(vars["id"]), "Owner": user.Id.Hex()}).One(&result)
	if err != nil {
		panic(jsonError{fmt.Sprintf("Could not find todo with id %s", vars["id"]), QueryError})
	}

	Json(rw, result)
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

func Login(rw http.ResponseWriter, req *http.Request) {
	session, err := MongoSession(req)
	if err != nil {
		panic(err)
	}

	user := User{}
	err = session.DB("todos").C("user").Find(bson.M{
		"Username": req.FormValue("Username"),
		"Password": req.FormValue("Password")}).One(&user)

	if err != nil {
		panic(jsonError{"Authentication failed", AuthenticationError})
	}

	fmt.Printf("%v", user)

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
