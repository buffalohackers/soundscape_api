package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type ArtistInfoResponse struct {
	Artist ArtistInfo `json:"artist"`
}

type ArtistInfo struct {
	Tags Tags `json:"tags"`
}

type Tags struct {
	Tag []TagInfo `json:"tag"`
}

type TagInfo struct {
	Name string `json:"name"`
}

const lastFmBaseUrl = "http://ws.audioscrobbler.com/2.0/?%s"

func GetArtistGenre(artist string) (string, error) {
	params := make(url.Values)
	params["method"] = []string{"artist.getinfo"}
	params["artist"] = []string{artist}
	params["api_key"] = []string{os.Getenv("LASTFMTOKEN")}
	params["format"] = []string{"json"}
	log.Println("GET:", fmt.Sprintf(lastFmBaseUrl, params.Encode()))
	resp, err := http.Get(fmt.Sprintf(lastFmBaseUrl, params.Encode()))
	if err != nil {
		log.Println("Could not GET Last.FM:", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read resp:", err.Error())
		return "", err
	}
	artistInfo := ArtistInfoResponse{}
	err = json.Unmarshal(body, &artistInfo)
	if err != nil {
		log.Println("Could not unmarshal JSON:", err.Error())
		return "", err
	}
	if len(artistInfo.Artist.Tags.Tag) == 0 {
		return "", errors.New("Artist not found")
	}
	return artistInfo.Artist.Tags.Tag[0].Name, nil
}
