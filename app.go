package main

import (
	"github.com/nickdirienzo/go-json-rest"
	"labix.org/v2/mongo"
)

type api struct {
	db *Session
}

func main() {

	session, err := mgo.Dial(url)

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
