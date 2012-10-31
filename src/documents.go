package main

import (
	"labix.org/v2/mgo/bson"
	"time"
)

const (
	StatusCompleted int = iota
	StatusActive
)

type Todo struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Owner     bson.ObjectId `bson:"ow" json:"owner"`
	Title     string        `bson:"ti" json:"title"`
	Priority  int           `bson:"pr" json:"priority"`
	Status    int           `bson:"st" json:"status"`
	DueAt     time.Time     `bson:"du" json:"due_at"`
	CreatedAt time.Time     `bson:"cr" json:"created_at"`
}

type User struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"u" json:"username"`
	Password string        `bson:"p" json:"password"`
}
