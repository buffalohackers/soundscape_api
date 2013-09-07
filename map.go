package main

import (
	"fmt"
	"github.com/nickdirienzo/go-json-rest"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://maps.googleapis.com/maps/api/place/autocomplete/json?" +
	"input=%s&types=(cities)&location=%s&sensor=true&key=AIzaSyAYfeOIDGoPo1A5_Wgm8J1MoMC2KdAuJBM"

func (self *Api) GMapsMirror(w *rest.ResponseWriter, r *rest.Request) {
	q := r.URL.Query()
	input := q.Get("input")
	location := q.Get("location")
	resp, _ := http.Get(fmt.Sprintf(baseUrl, input, location))
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
}
