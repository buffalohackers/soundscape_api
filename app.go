package main

import (
	"github.com/nickdirienzo/go-json-rest"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

type Api struct {
	DbName       string
	MongoSession *mgo.Session
}

func (self *Api) GetSongs(w rest.ResponseWriter, r *rest.Request) {

}

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Could not get Mongo")
	}
	api := Api{MongoSession: session, DbName: "mugo"}

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.RouteObjectMethod("GET", "/sessions", &api, "GetSessions"),
		rest.RouteObjectMethod("POST", "/songs", &api, "PostSongs"),
		rest.RouteObjectMethod("GET", "/songs", &api, "GetSongs"),
	)

	hostname, port := "127.0.0.1", "8080"
	log.Println("Starting server on " + hostname + ":" + port)
	http.ListenAndServe(hostname+":"+port, &handler)
}
