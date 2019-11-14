package twitch

import (
	"net/http"
	"os/exec"
	"time"
)

type Client struct {
	user          *UserApi
	streamHandler LiveStream
}

const clientId = "g0m4aoe1qgv0lqais31yp27yzvw603"

func NewClient() *Client {
	http := &http.Client{
		Timeout: time.Second * 2,
	}

	accessToken := GetAccessToken()

	return &Client{
		user: &UserApi{
			ClientId:    clientId,
			AccessToken: accessToken,
			Http:        http,
		},
		streamHandler: &HttpLiveStream{
			provider: &TwitchProvider{
				ClientId: clientId,
				Http:     http,
			},
		},
	}
}

func (client Client) Play(stream string) {
	vlc := exec.Command("cvlc", "-")
	pipe, _ := vlc.StdinPipe()

	client.streamHandler.Start(stream, pipe)

	vlc.Run()
}

func (client Client) ListStreams() []Stream {
	return client.user.GetLiveStreams()
}
