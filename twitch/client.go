package twitch

import (
	"net/http"
	"time"
)

type Client struct {
	User   *UserApi
	Player *Player
}

func NewClient(clientId string, accessToken string) *Client {
	http := &http.Client{
		Timeout: time.Second * 2,
	}

	return &Client{
		User: &UserApi{
			ClientId:    clientId,
			AccessToken: accessToken,
			Http:        http,
		},
		Player: &Player{
			Stream: &StreamApi{
				ClientId: clientId,
				Http:     http,
			},
		},
	}
}
