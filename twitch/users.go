package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"display_name"`
	Login string `json:"login"`
	Image string `json:"profile_image_url"`
}

type Stream struct {
	Channel Channel `json:"channel"`
}

type Channel struct {
	Name string `json:"name"`
	Game string `json:"game"`
}

type UserApi struct {
	ClientId    string
	AccessToken string
	Http        *http.Client
}

const USER_URL = "https://api.twitch.tv/helix/users?login=%s"
const LIVE_CHANNELS_URL = "https://api.twitch.tv/kraken/streams/followed?stream_type=live"

func (api UserApi) GetUser(user string) User {
	response := api.get(fmt.Sprintf(USER_URL, user))
	body := struct{ Data []User }{}
	err := json.NewDecoder(response.Body).Decode(&body)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return body.Data[0]
}

func (api UserApi) GetLiveStreams() []Stream {
	request, _ := http.NewRequest("GET", LIVE_CHANNELS_URL, nil)
	request.Header.Add("Client-ID", api.ClientId)
	request.Header.Add("Authorization", fmt.Sprintf("OAuth %s", api.AccessToken))
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

func (api UserApi) get(url string) *http.Response {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Client-ID", api.ClientId)
	response, err := api.Http.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	return response
}
