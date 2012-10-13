package todos

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Todo struct {
	Id        bson.ObjectId "_id,omitempty"
	Owner     bson.ObjectId
	Title     string
	Priority  int
	Status    int
	DueAt     time.Time
	CreatedAt time.Time
}

type User struct {
	Id       bson.ObjectId "_id,omitempty"
	Username string
	Password string
}
