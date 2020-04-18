package model

import "gopkg.in/mgo.v2/bson"

type MainPage struct {
	Id            bson.ObjectId `bson:"-"`
	Caption       string        `bson:"caption" description:"说说内容"`
	Code          string        `bson:"code" description:"数据库code"`
	FullName      string        `bson:"full_name" description:"用户全名"`
	Username      string        `bson:"username" description:"跳转用户主页"`
	Data          string        `bson:"time" description:"发表日期"`
	ThumbnailSrc  string        `bson:"thumbnailsrc" description:"数据库code"`
	MediaType	  float64		`bson:"mediaType" description:"类型"`
	ProfilePicUrl string        `bson:"profilepicurl" description:"用户头像"`
	From 		  string 		`bson:"from" description:"标识进来的路径"`
	InfoPage      []MainPage    `bson:"infopage" description:"用户信息"`
}