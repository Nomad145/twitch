package twitch

import (
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"log"
	"net/http"
)

type AccessToken struct {
	Sig   string
	Token string
}

type StreamPlaylist struct {
	Playlist m3u8.Playlist
	ListType m3u8.ListType
}

type Stream struct {
	Quality  string
	Playlist m3u8.Playlist
}

type StreamApi struct {
	ClientId string
	Http     *http.Client
}

const TOKEN_URL = "https://api.twitch.tv/channels/%s/access_token"

const HLS_URL = "https://usher.ttvnw.net/api/channel/hls/%s.m3u8"

func (api StreamApi) GetStreamPlaylist(user User) StreamPlaylist {
	token := api.getAccessToken(user)
	playlist := api.getHLS(user, token)

	return playlist
}

func (api StreamApi) getAccessToken(user User) AccessToken {
	response := api.get(fmt.Sprintf(TOKEN_URL, user.Login), map[string]string{})
	token := AccessToken{}
	err := json.NewDecoder(response.Body).Decode(&token)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return token
}

func (api StreamApi) get(url string, params map[string]string) *http.Response {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Client-ID", api.ClientId)
	query := request.URL.Query()

	for key, value := range params {
		query.Add(key, value)
	}

	request.URL.RawQuery = query.Encode()

	response, err := api.Http.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	return response
}

func (api StreamApi) getHLS(user User, token AccessToken) StreamPlaylist {
	params := map[string]string{
		"sig":              token.Sig,
		"token":            token.Token,
		"player":           "twitchweb",
		"p":                "707780",
		"type":             "any",
		"allow_source":     "true",
		"allow_audio_only": "true",
		"allow_spectre":    "false",
	}

	response := api.get(fmt.Sprintf(HLS_URL, user.Login), params)

	playlist, listtype, _ := m3u8.DecodeFrom(response.Body, true)

	return StreamPlaylist{
		playlist,
		listtype,
	}
}

func (api StreamApi) GetBest(playlist StreamPlaylist) Stream {
	masterPlaylist := playlist.Playlist.(*m3u8.MasterPlaylist)
	variant := masterPlaylist.Variants[0]

	response := api.get(variant.URI, map[string]string{})

	mediaPlaylist, _, _ := m3u8.DecodeFrom(response.Body, true)

	return Stream{
		variant.VariantParams.Resolution,
		mediaPlaylist,
	}
}
