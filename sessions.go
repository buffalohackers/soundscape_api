package main

import (
	"github.com/nickdirienzo/go-json-rest"
	// "labix.org/v2/mgo"
	"crypto/md5"
	"io"
	"log"
	"net/http"
	"time"
	"fmt"
)

type session struct {
	SessionKey  string   `bson:"session_key"`
	SongsPlayed []string `bson:"songs_played"`
}

func (self *Api) GetSessions(w *rest.ResponseWriter, r *rest.Request) {
	d := r.RemoteAddr + time.Now().String()
	h := md5.New()
	io.WriteString(h, d)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	session := session{SessionKey: hash}
	k := self.Db.DB(self.DbName).C("session-keys")
	err := k.Insert(&session)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError, "get.sessions")
	}

	response := Response{}
	response["data"] = "Bullshit"
	w.WriteJson(&response, http.StatusOK)
}
