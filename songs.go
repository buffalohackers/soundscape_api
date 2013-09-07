package main

import (
	"github.com/nickdirienzo/go-json-rest"
	"labix.org/v2/mgo/bson"
	"log"
	"math"
	"net/http"
	"sort"
)

type Song struct {
	Id    string  `bson:"id" json:"id"`
	Lat   float64 `bson:"lat" json:"lat"`
	Long  float64 `bson:"long" json:"long"`
	Genre string  `bson:"genre" json:"genre"`
}

type SortedSong struct {
	Id   string  `bson:"id" json:"id"`
	Dist float64 `bson:"dist" json:"dist"`
}

type SortedSongs []SortedSong

func (s SortedSongs) Len() int {
	return len(s)
}

func (s SortedSongs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedSongs) Less(i, j int) bool {
	return s[i].Dist < s[j].Dist
}

type ClosestSongQuery struct {
	SessionKey string  `bson:"session_key" json:"session_key"`
	Lat        float64 `bson:"lat" json:"lat"`
	Long       float64 `bson:"long" json:"long"`
	Genre      string  `bson:"genre" json:"genre"`
}

type SongResponse struct {
	Success bool `bson:"success" json:"success"`
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

func calculateDistances(query ClosestSongQuery, songs []Song) SortedSongs {
	var sortedSongs SortedSongs
	for _, s := range songs {
		dist := math.Sqrt(math.Pow(query.Lat-s.Lat, 2) + math.Pow(query.Long-s.Long, 2))
		sortedSong := SortedSong{Id: s.Id, Dist: dist}
		sortedSongs = append(sortedSongs, sortedSong)
	}
	return sortedSongs
}

func (self *Api) getUnlistenedSong(query ClosestSongQuery, songs SortedSongs) SortedSong {
	var session Session
	sessions := self.MongoSession.DB(self.DbName).C("sessions")
	sessions.Find(bson.M{"session_key": query.SessionKey}).One(&session)
	for _, song := range songs {
		if !session.SongsPlayed[song.Id] {
			return song
		}
	}
	return songs[0]
}

func (self *Api) getClosestSong(query ClosestSongQuery) SortedSong {
	songs := []Song{}
	songCollection := self.MongoSession.DB(self.DbName).C("songs")
	if query.Genre == "" {
		songCollection.Find(nil).All(&songs)
	} else {
		songCollection.Find(bson.M{"genre": query.Genre}).All(&songs)
	}
	sortedSongs := calculateDistances(query, songs)
	sort.Sort(sortedSongs)
	sortedSong := self.getUnlistenedSong(query, sortedSongs)
	return sortedSong
}

func (self *Api) GetSongs(w *rest.ResponseWriter, r *rest.Request) {
	method := "songs.get"
	req := ClosestSongQuery{}
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		log.Println("Could not decode request:", err.Error())
		rest.Error(w, "Could not process songs.get request", http.StatusBadRequest, method)
		return
	}
	closestSong := self.getClosestSong(req)
	_ = self.updateSessionSongs(req, closestSong)
	w.WriteJson(&closestSong, http.StatusOK)
}
