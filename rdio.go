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
	hackersBaseUrl = "http://localhost:8080/%s"
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
