package main

import (
	"fmt"
	"github.com/nickdirienzo/go-json-rest"
	"github.com/nickdirienzo/rdigo"
	"log"
	"net/http"
	"os"
	// "time"
)

const (
	month          = 60 * 60 * 24 * 30
	hackersBaseUrl = "http://buffalohackers.com/%s"
)

type Redirect struct {
	See string `json:"see"`
}

func getCookie(r *rest.Request, name string) *http.Cookie {
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func newRdioClient() rdigo.Rdio {
	return rdigo.NewClient(os.Getenv("RDIOTOKEN"), os.Getenv("RDIOSECRET"))
}

func authedRdioClient(at, ats string) rdigo.Rdio {
	return rdigo.AuthenticatedClient(os.Getenv("RDIOTOKEN"), os.Getenv("RDIOSECRET"), at, ats)
}

func (self *Api) LogIn(w *rest.ResponseWriter, r *rest.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "at",
		Value: "",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "ats",
		Value: "",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "rt",
		Value: "",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "rts",
		Value: "",
	})
	rdio := newRdioClient()
	rToken, url, err := rdio.BeginAuthentication(fmt.Sprintf(hackersBaseUrl, "rdio"))
	if err != nil {
		log.Println("Rdio Auth Error:", err.Error())
		rest.Error(w, "Could not authenticate with Rdio.", http.StatusBadRequest, "login.get")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "rt",
		Value: rToken.Token,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "rts",
		Value: rToken.Secret,
	})
	w.WriteJson(&Redirect{See: url}, http.StatusTemporaryRedirect)
}

func (self *Api) RdioCallback(w *rest.ResponseWriter, r *rest.Request) {
	rt := getCookie(r, "rt")
	rts := getCookie(r, "rts")
	if rt == nil || rts == nil {
		log.Println("rt or rts not set")
		rest.Error(w, "Could not authenticate with Rdio.", http.StatusBadRequest, "rdio.get")
		return
	}
	verifier := r.URL.Query().Get("oauth_verifier")
	log.Println(verifier)
	rdio := newRdioClient()
	err := rdio.CompleteAuthentication(rt.Value, rts.Value, verifier)
	if err != nil {
		log.Println("Rdio Auth Error:", err.Error())
		rest.Error(w, "Could not authenticate with Rdio.", http.StatusBadRequest, "rdio.get")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "at",
		Value: rdio.AccessToken.Token,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "ats",
		Value: rdio.AccessToken.Secret,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "rt",
		Value: "",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "rts",
		Value: "",
	})
	w.WriteJson(&Redirect{See: fmt.Sprintf(hackersBaseUrl, "gpsongs/")}, http.StatusTemporaryRedirect)
}

func (self *Api) GetPlaybackToken(w *rest.ResponseWriter, r *rest.Request) {
	at := getCookie(r, "at")
	ats := getCookie(r, "ats")
	if at == nil || ats == nil {
		log.Println("at or ats not found.")
		rest.Error(w, "Could not authenticate with Rdio.", http.StatusBadRequest, "playbackToken.get")
		return
	}
	rdio := authedRdioClient(at.Value, ats.Value)
	ret, err := rdio.GetPlaybackToken("buffalohackers.com")
	if err != nil {
		log.Println("Rdio Call Fail:", err.Error())
		rest.Error(w, "Rdio Call Failed", http.StatusBadRequest, "playbackToken.get")
		return
	}
	w.WriteJson(&ret, http.StatusOK)
}

func (self *Api) SearchRdio(w *rest.ResponseWriter, r *rest.Request) {
	at := getCookie(r, "at")
	ats := getCookie(r, "ats")
	if at == nil || ats == nil {
		log.Println("at or ats not found.")
		rest.Error(w, "Could not authenticate with Rdio.", http.StatusBadRequest, "search.get")
		return
	}
	q := r.URL.Query().Get("q")
	rdio := authedRdioClient(at.Value, ats.Value)
	query := make(map[string]string)
	query["query"] = q
	query["types"] = "Artist,Album,Track"
	ret, err := rdio.Call("search", query)
	if err != nil {
		log.Println("Rdio Call Fail:", err.Error())
		rest.Error(w, "Rdio Call Failed", http.StatusBadRequest, "search.get")
		return
	}
	w.WriteJson(&ret, http.StatusOK)
}
