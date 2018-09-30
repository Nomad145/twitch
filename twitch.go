package twitch

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Id    int    `json:"id,string"`
	Name  string `json:"display_name"`
	Image string `json:"profile_image_url"`
}

type Client struct {
	ClientID string
	Http     *http.Client
}

func NewClient(client_id string) Client {
	return &Client{
		ClientID: client_id,
		Http: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (client Client) GetUser(user string) User {
	response := get("https://api.twitch.tv/helix/users?login=" + user)

	body := struct{ Data []User }{}

	err = json.NewDecoder(res.Body).Decode(&body)

	if err != nil {
		log.Fatal(err)
	}

	return body.Data[0]
}

func (client Client) get(url string) http.Response {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Client-ID", client.ClientID)
	response, err := client.Do(request)

	if err != nil {
		log.fatal(err)
	}

	defer response.Body.Close()

	return response
}
