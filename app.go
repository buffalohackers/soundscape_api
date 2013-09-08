package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/nickdirienzo/go-json-rest"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

type Api struct {
	DbName       string
	MongoSession *mgo.Session
}

func StartWebSocketServer() {
	http.Handle("/ws", websocket.Handler(wsHandler))
	hostname, port := "127.0.0.1", "8081"
	log.Println("Starting WS Server on " + hostname + ":" + port)
	err := http.ListenAndServe(hostname+":"+port, nil)
	if err != nil {
		log.Fatal("Could not start WS server:", err.Error())
	}
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
		rest.RouteObjectMethod("GET", "/map", &api, "GMapsMirror"),
		rest.RouteObjectMethod("GET", "/mockPins", &api, "GenerateMockPins"),
		rest.RouteObjectMethod("GET", "/login", &api, "LogIn"),
		rest.RouteObjectMethod("GET", "/rdio", &api, "RdioCallback"),
		rest.RouteObjectMethod("GET", "/playbackToken", &api, "GetPlaybackToken"),
		rest.RouteObjectMethod("GET", "/search", &api, "SearchRdio"),
		rest.RouteObjectMethod("GET", "/allSongs", &api, "GetAllSongs"),
	)
	go h.run()
	go StartWebSocketServer()
	hostname, port := "127.0.0.1", "8080"
	log.Println("Starting server on " + hostname + ":" + port)
	http.ListenAndServe(hostname+":"+port, &handler)
}
