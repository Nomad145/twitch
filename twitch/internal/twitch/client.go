package twitch

import (
	"fmt"
	"github.com/lukechampine/hey"
	"net/http"
	"os/exec"
	"time"
)

type Client struct {
	user          *UserApi
	streamHandler LiveStream
}

const publicClientId = "g0m4aoe1qgv0lqais31yp27yzvw603"
const secretClientId = "kimne78kx3ncx6brgo4mv6wki5h1ko"

func NewClient() *Client {
	http := &http.Client{
		Timeout: time.Second * 2,
	}

	accessToken := GetAccessToken()

	return &Client{
		user: &UserApi{
			ClientId:    publicClientId,
			AccessToken: accessToken,
			Http:        http,
		},
		streamHandler: &HttpLiveStream{
			provider: &TwitchProvider{
				ClientId: secretClientId,
				Http:     http,
			},
		},
	}
}

func (client Client) Play(stream string) {
	vlc := exec.Command("cvlc", "-")
	pipe, _ := vlc.StdinPipe()

	hey.Push(hey.Notification{
		Title:    "Twitch",
		Body:     fmt.Sprintf("Starting %s's stream shortly...", stream),
		AppName:  "twitch",
		Duration: time.Second * 15,
	})

	client.streamHandler.Start(stream, pipe)

	vlc.Run()
}

func (client Client) ListStreams() []Stream {
	return client.user.GetLiveStreams()
}
