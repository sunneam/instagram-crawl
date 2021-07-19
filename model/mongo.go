package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var(
	Db *mgo.Database
	url =""
)
func Init()  {

	session,err:=mgo.Dial(url)
	if err!=nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Monotonic, true)
	Db=session.DB("veryins")
}