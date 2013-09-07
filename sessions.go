package main

import (
	"crypto/md5"
	"fmt"
	"github.com/nickdirienzo/go-json-rest"
	"io"
	"log"
	"net/http"
	"time"
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
	k := self.MongoSession.DB(self.DbName).C("sessions")
	err := k.Insert(&session)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError, "sessions.get")
	}

	response := Response{}
	response["data"] = "Bullshit"
	w.WriteJson(&response, http.StatusOK)
}
