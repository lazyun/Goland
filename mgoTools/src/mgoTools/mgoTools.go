package mgoTools

import (
	"gopkg.in/mgo.v2"
)


type MgoTools struct {
	MongoUri	string
	MongoPass	string
	MongoUriStr string

	mgoSession 	*mgo.Session

}


func (this *MgoTools) SetConfig(uri , passswd string) bool {
	mgo.Dial()
}