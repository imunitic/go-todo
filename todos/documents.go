package todos

import (
	"labix.org/v2/mgo/bson"
	"time"
)

const (
	StatusCompleted int = 0
	StatusActive    int = 1
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
