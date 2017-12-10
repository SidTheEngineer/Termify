package auth

import (
	"Termify/helpers"
	"fmt"
	"net/http"
	"os"

	"github.com/skratchdot/open-golang/open"
)

const (
	authURL         = "https://accounts.spotify.com/authorize"
	authRedirectURI = "http://localhost:8000/callback"
	scopes          = "user-read-playback-state"
)

// Authorize sends a request to Spotify's authorize URL with the
// appropriate headers and query parameters. A new browser window
// will be opened and the user will be prompted to login and or
// provide Termify access to their Spotify account, after which,
// the user is redirected to /callback.
func Authorize() {
	spotifyAuthURL := createSpotifyAuthURL()
	open.Run(spotifyAuthURL)
}

// GetClientIDAndSecret returns the Spotify client ID and client secret
// from environment variables on the machine the application is run on.
func GetClientIDAndSecret() (string, string) {
	return os.Getenv("SPOTIFY_CLIENT"), os.Getenv("SPOTIFY_SECRET")
}

func createSpotifyAuthURL() string {
	req, err := http.NewRequest("GET", authURL, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addAuthQueryParams(req)

	return req.URL.String()
}

func addAuthQueryParams(req *http.Request) {
	clientID, _ := GetClientIDAndSecret()
	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", authRedirectURI)
	q.Add("state", helpers.GenerateRandomString(32))
	q.Add("scope", scopes)
	req.URL.RawQuery = q.Encode()
}
