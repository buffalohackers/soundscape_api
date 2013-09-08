package main

import (
	"fmt"
	"github.com/nickdirienzo/go-json-rest"
	"labix.org/v2/mgo/bson"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

type Song struct {
	SessionKey string  `bson:"session_key" json:"session_key"`
	Url        string  `bson:"url" json:"url"`
	Artist     string  `bson:"artist" json:"artist"`
	SongName   string  `bson:"song" json:"song"`
	Id         string  `bson:"id" json:"id"`
	Lat        float64 `bson:"lat" json:"lat"`
	Long       float64 `bson:"long" json:"long"`
	Genre      string  `bson:"genre" json:"genre"`
}

func (s *Song) String() string {
	return fmt.Sprintf("{\"id\": \"%s\", \"lat\": %s, \"long\": %s, \"genre\": \"%s\", \"url\": \"%s\", \"artist\": \"%s\", \"song\": \"%s\"}",
		s.Id, strconv.FormatFloat(s.Lat, 'G', -1, 64), strconv.FormatFloat(s.Long, 'G', -1, 64), s.Genre, s.Url, s.Artist, s.SongName)
}

type SortedSong struct {
	Id   string  `bson:"id" json:"id"`
	Dist float64 `bson:"dist" json:"dist"`
	Lat  float64 `bson:"lat" json:"lat"`
	Long float64 `bson:"long" json:"long"`
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
	genre, err := GetArtistGenre(song.Artist)
	if err != nil {
		log.Println("Could not get genre:", err.Error())
		rest.Error(w, fmt.Sprintf("Could not get genre for artist %s", song.Artist), http.StatusBadRequest, method)
		return
	}
	song.Genre = genre
	songCollection := self.MongoSession.DB(self.DbName).C("songs")
	err = songCollection.Insert(&song)
	if err != nil {
		log.Println("Could not insert song:", err.Error())
		rest.Error(w, "Could not process song", http.StatusInternalServerError, method)
		return
	}
	err = self.updateSessionSongs(song.SessionKey, song.Id)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, fmt.Sprintf("Could not update session songs:%s", err.Error()), http.StatusBadRequest, method)
		return
	}
	resp := SongResponse{Success: true}
	h.broadcast <- song.String()
	w.WriteJson(&resp, http.StatusCreated)
}

func calculateDistances(query ClosestSongQuery, songs []Song) SortedSongs {
	var sortedSongs SortedSongs
	for _, s := range songs {
		dist := math.Sqrt(math.Pow(query.Lat-s.Lat, 2) + math.Pow(query.Long-s.Long, 2))
		sortedSong := SortedSong{Id: s.Id, Dist: dist, Long: s.Long, Lat: s.Lat}
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
	self.clearSessionSongs(query.SessionKey)
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
	if len(sortedSongs) > 0 {
		return sortedSongs[0]
	} else {
		return nil
	}
	//sortedSong := self.getUnlistenedSong(query, sortedSongs)
	// return sortedSong
}

func (self *Api) GetSongs(w *rest.ResponseWriter, r *rest.Request) {
	method := "songs.get"
	query := r.URL.Query()
	sessionKey := query.Get("session_key")
	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		log.Println("Cannot parse lat:", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest, method)
		return
	}
	long, err := strconv.ParseFloat(query.Get("long"), 64)
	if err != nil {
		log.Println("Cannot parse long:", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest, method)
		return
	}
	genre := query.Get("genre")
	req := ClosestSongQuery{SessionKey: sessionKey, Lat: lat, Long: long, Genre: genre}
	closestSong := self.getClosestSong(req)
	err = self.updateSessionSongs(req.SessionKey, closestSong.Id)
	if err != nil {
		log.Println("Cannot update session songs:", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest, method)
		return
	}
	w.WriteJson(&closestSong, http.StatusOK)
}

func (self *Api) GetAllSongs(w *rest.ResponseWriter, r *rest.Request) {
	songCollection := self.MongoSession.DB(self.DbName).C("songs")
	var songs []Song
	_ = songCollection.Find(nil).All(&songs)
	w.WriteJson(&songs, http.StatusOK)
}
