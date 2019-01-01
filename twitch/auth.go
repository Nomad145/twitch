package twitch

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
	"log"
	"os"
)

func GetAccessToken() string {
	if tokenExists() {
		token := getStoredToken()
		token = refreshToken(token)

		return token.AccessToken
	}

	token := getNewToken()
	storeToken(token)

	return token.AccessToken
}

func buildConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		Scopes:       []string{"user:read:email"},
		RedirectURL:  "http://localhost",
		Endpoint:     twitch.Endpoint,
	}
}

func getNewToken() oauth2.Token {
	context := context.Background()
	config := buildConfig()
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := config.Exchange(context, code)
	if err != nil {
		log.Fatal(err)
	}

	return *tok
}

func tokenExists() bool {
	_, err := os.Stat("/home/nomad/.config/twitch")

	return err == nil
}

func getStoredToken() oauth2.Token {
	var token oauth2.Token

	file, err := os.Open("/home/nomad/.config/twitch")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	json.NewDecoder(file).Decode(&token)

	return token
}

func refreshToken(token oauth2.Token) oauth2.Token {
	context := context.Background()
	tokenSource := buildConfig().TokenSource(context, &token)
	refreshedToken, err := tokenSource.Token()

	if err != nil {
		log.Fatal(err)
	}

	storeToken(*refreshedToken)

	return *refreshedToken
}

func storeToken(token oauth2.Token) {
	file, err := os.Create("/home/nomad/.config/twitch")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	json.NewEncoder(file).Encode(token)
}
