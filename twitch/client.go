package twitch

import (
	"net/http"
	"time"
)

type Client struct {
	User   *UserApi
	Player *Player
}

func NewClient() *Client {
	http := &http.Client{
		Timeout: time.Second * 2,
	}

	accessToken := GetAccessToken()
	clientId := "g0m4aoe1qgv0lqais31yp27yzvw603"

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
