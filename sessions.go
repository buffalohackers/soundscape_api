package main

import (
	"crypto/md5"
	"fmt"
	"github.com/nickdirienzo/go-json-rest"
	"io"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"time"
)

type Session struct {
	SessionKey  string          `bson:"session_key" json:"session_key"`
	SongsPlayed map[string]bool `bson:"songs_played" json:"songs_played"`
}

func (self *Api) GetSessions(w *rest.ResponseWriter, r *rest.Request) {
	d := r.RemoteAddr + time.Now().String()
	h := md5.New()
	io.WriteString(h, d)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	session := Session{SessionKey: hash}
	k := self.MongoSession.DB(self.DbName).C("sessions")
	err := k.Insert(&session)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError, "sessions.get")
	}
	w.WriteJson(&session, http.StatusOK)
}

func (self *Api) updateSessionSongs(sessionKey, songId string) error {
	var s Session
	sessions := self.MongoSession.DB(self.DbName).C("sessions")
	sessions.Find(bson.M{"session_key": sessionKey}).One(&s)
	s.SongsPlayed[songId] = true
	err := sessions.Update(bson.M{"session_key": sessionKey}, &s)
	return err
}

func (self *Api) clearSessionSongs(sessionKey string) error {
	var s Session
	sessions := self.MongoSession.DB(self.DbName).C("sessions")
	sessions.Find(bson.M{"session_key": sessionKey}).One(&s)
	s.SongsPlayed = make(map[string]bool)
	err := sessions.Update(bson.M{"session_key": sessionKey}, &s)
	return err
}
