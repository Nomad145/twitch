package twitch

import (
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AccessToken struct {
	Sig   string
	Token string
}

type StreamPlaylist struct {
	Playlist m3u8.Playlist
	ListType m3u8.ListType
	URL      string
}

type StreamApi struct {
	ClientId string
	Http     *http.Client
}

const TOKEN_URL = "https://api.twitch.tv/channels/%s/access_token"
const HLS_URL = "https://usher.ttvnw.net/api/channel/hls/%s.m3u8"

func (api StreamApi) GetMasterPlaylist(user string) StreamPlaylist {
	token := api.getAccessToken(user)
	playlist := api.fetchMasterPlaylist(user, token)

	log.Print("Starting Stream...")

	return playlist
}

func (api StreamApi) getAccessToken(user string) AccessToken {
	token := AccessToken{}
	request, _ := http.NewRequest("GET", fmt.Sprintf(TOKEN_URL, user), nil)
	request.Header.Add("Client-ID", api.ClientId)

	response, err := api.Http.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(response.Body).Decode(&token)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return token
}

func (api StreamApi) fetchMasterPlaylist(user string, token AccessToken) StreamPlaylist {
	request, _ := http.NewRequest("GET", fmt.Sprintf(HLS_URL, user), nil)

	params := url.Values{}
	params.Set("sig", token.Sig)
	params.Set("token", token.Token)
	params.Set("player", "twitchweb")
	params.Set("p", "707790")
	params.Set("type", "any")
	params.Set("allow_source", "true")
	params.Set("allow_audio_only", "true")
	params.Set("allow_spectre", "false")
	request.URL.RawQuery = params.Encode()

	response, err := api.Http.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	playlist, listType, _ := m3u8.DecodeFrom(response.Body, true)

	response.Body.Close()

	return StreamPlaylist{
		playlist,
		listType,
		request.URL.String(),
	}
}

func (api StreamApi) GetMediaPlaylist(playlist StreamPlaylist) StreamPlaylist {
	masterPlaylist := playlist.Playlist.(*m3u8.MasterPlaylist)
	variant := masterPlaylist.Variants[0]

	response, err := http.Get(variant.URI)

	if err != nil {
		log.Fatal(err)
	}

	mediaPlaylist, listType, _ := m3u8.DecodeFrom(response.Body, true)

	return StreamPlaylist{
		mediaPlaylist,
		listType,
		variant.URI,
	}
}

func (api StreamApi) RefreshPlaylist(playlist StreamPlaylist) StreamPlaylist {
	response, _ := http.Get(playlist.URL)
	refreshedPlaylist, listType, err := m3u8.DecodeFrom(response.Body, true)

	if refreshedPlaylist == nil {
		log.Println(err)
		log.Println("Playlist was nil...")
	}

	return StreamPlaylist{
		refreshedPlaylist,
		listType,
		playlist.URL,
	}
}

func (api StreamApi) DownloadSegment(url string, target io.Writer) {
	segment, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(segment.Body)

	if err != nil {
		log.Fatal(err)
	}

	target.Write(body)
	segment.Body.Close()
}
