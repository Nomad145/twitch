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

const clientId = "kimne78kx3ncx6brgo4mv6wki5h1ko"

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
