package twitch

import (
	"net/http"
	"time"
)

type Client struct {
	User     *UserApi
	Stream   *StreamApi
}

func NewClient(clientId string) *Client {
	http := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Client{
		User: &UserApi{
			ClientId: clientId,
			Http: http,
		},
		Stream: &StreamApi{
			ClientId: clientId,
			Http: http,
		},
	}
}
