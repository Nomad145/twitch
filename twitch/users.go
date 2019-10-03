package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Stream struct {
	Channel Channel `json:"channel"`
}

type Channel struct {
	Name   string `json:"name"`
	Game   string `json:"game"`
	Status string `json:"status"`
}

type UserApi struct {
	ClientId    string
	AccessToken string
	Http        *http.Client
}

const LIVE_CHANNELS_URL = "https://api.twitch.tv/kraken/streams/followed?stream_type=live"

func (api UserApi) GetLiveStreams() []Stream {
	request, _ := http.NewRequest("GET", LIVE_CHANNELS_URL, nil)
	request.Header.Add("Client-ID", api.ClientId)
	request.Header.Add("Authorization", fmt.Sprintf("OAuth %s", api.AccessToken))
	request.Header.Add("Accept", "application.vnd.twitchtv.v5+json")
	response, err := api.Http.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	body := struct {
		Streams []Stream `json:"streams"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&body)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return body.Streams
}
