package twitch

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"display_name"`
	Login string `json:"login"`
	Image string `json:"profile_image_url"`
}

type UserApi struct {
	ClientId string
	Http *http.Client
}

func (api UserApi) Get(user string) User {
	response := api.get("https://api.twitch.tv/helix/users?login=" + user)
	body := struct{ Data []User }{}
	err := json.NewDecoder(response.Body).Decode(&body)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return body.Data[0]
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
