package main

import (
	"crypto/md5"
	"github.com/nickdirienzo/go-rest-json"
	"labix.org/v2/mgo"
	"net/http"
)

type session struct {
	SessionKey  string
	SongsPlayed []string
}

func (self *Api) GetSession(w *rest.ResponseWriter, r *rest.Request) {
	ip := r.Headers.Get("Remote_Addr")
	log.Println(ip)

	session := session{}
	k := self.Db(self.DbName).C("session-keys")
	err := k.Insert(&session)

	return fmt.Fprintln("FUCK")
}
