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
	fmt.Println("测试")
	session.SetMode(mgo.Monotonic, true)
	Db=session.DB("veryins")
	fmt.Println("连接数据库成功")
}