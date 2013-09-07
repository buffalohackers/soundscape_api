package main

import (
	"github.com/nickdirienzo/go-json-rest"
	"log"
	"net/http"
)

type Song struct {
	Id    string  `bson:"id"`
	Lat   float64 `bson:"lat"`
	Long  float64 `bson:"long"`
	Genre string  `bson:"genre"`
}

type SongResponse struct {
	Success bool `bson:"success"`
}

func (self *Api) PostSongs(w *rest.ResponseWriter, r *rest.Request) {
	method := "songs.post"
	song := Song{}
	err := r.DecodeJsonPayload(&song)
	if err != nil {
		log.Println("Could not decode song:", err.Error())
		rest.Error(w, "Could not process song", http.StatusBadRequest, method)
		return
	}
	songCollection := self.MongoSession.DB(self.DbName).C("songs")
	err = songCollection.Insert(&song)
	if err != nil {
		log.Println("Could not insert song:", err.Error())
		rest.Error(w, "Could not process song", http.StatusInternalServerError, method)
		return
	}
	resp := SongResponse{Success: true}
	w.WriteJson(&resp, http.StatusCreated)
}
