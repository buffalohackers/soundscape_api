package main

import (
	"github.com/nickdirienzo/go-json-rest"
	"math/rand"
	"net/http"
	"strconv"
)

type MockPin struct {
	Lat   float64
	Long  float64
	Genre string
}

var genres = []string{"rock", "pop", "rap", "country", "dubstep", "electro house"}

func GenerateMockPin() MockPin {
	lat := (rand.Float64() * 25) + 24.52
	long := -1 * ((rand.Float64() * 57) + 66.95)
	genre := genres[rand.Intn(len(genres))]
	return MockPin{Lat: lat, Long: long, Genre: genre}
}

func (self *Api) GenerateMockPins(w *rest.ResponseWriter, r *rest.Request) {
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest, "mockPins.get")
		return
	}
	var pins []MockPin
	for i := 0; i < n; i++ {
		pins = append(pins, GenerateMockPin())
	}
	w.WriteJson(&pins, http.StatusOK)
}
