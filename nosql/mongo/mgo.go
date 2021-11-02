package mongo

import (
	"github.com/globalsign/mgo"
	"time"
	"fmt"
	"huashan/config"
)

var (
	MGO *mgo.Session
)

func Initialize()  {
	if MGO == nil {
		MGO = NewMgo()
	}
}

func UnInitialize() {
	if MGO != nil {
		MGO.Close()
	}
}

func NewMgo() *mgo.Session {
	if MGO != nil {
		return MGO
	}
	host := config.Get("mgo.host")
	username := config.Get("mgo.username")
	pwd := config.Get("mgo.password")
	database := config.Get("mgo.database")
	MGO,err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:[]string{host},
		Username:username,
		Password:pwd,
		Database:database,
		Timeout:30 * time.Second,
	})

	if err != nil {
		panic(fmt.Errorf("connect mongodb is error : %v",err))
	}
	return MGO
}



