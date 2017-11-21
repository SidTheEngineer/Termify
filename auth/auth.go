package auth

import (
	"Termify/helpers"
	"fmt"
	"net/http"
	"os"

	"github.com/skratchdot/open-golang/open"
)

// Authorize sends a request to Spotify's authorize URL with the
// appropriate headers and query parameters. A new browser window
// will be opened and the user will be prompted to login and or
// provide Termify access to their Spotify account.
func Authorize() {
	spotifyAuthURL := createSpotifyAuthURL()
	open.Run(spotifyAuthURL)
}

func createSpotifyAuthURL() string {
	req, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addAuthQueryParams(req)

	return req.URL.String()
}

func addAuthQueryParams(req *http.Request) {
	clientID, _ := getClientIDAndSecret()
	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", "http://localhost:8000/callback")
	q.Add("state", helpers.GenerateRandomString(32))
	req.URL.RawQuery = q.Encode()
}

func getClientIDAndSecret() (string, string) {
	return os.Getenv("SPOTIFY_CLIENT"), os.Getenv("SPOTIFY_SECRET")
}
