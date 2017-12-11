package auth

import (
	"Termify/util"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/skratchdot/open-golang/open"
)

const (
	authURL          = "https://accounts.spotify.com/authorize"
	authRedirectURI  = "http://localhost:8000/callback"
	scopes           = "user-read-playback-state user-modify-playback-state"
	tokenURL         = "https://accounts.spotify.com/api/token"
	tokenGrantType   = "authorization_code"
	tokenRedirectURI = "http://localhost:8000/callback"
)

// AccessToken is a token response returned after a Spotify authorization
// request is made.
type AccessToken struct {
	Token        string `json:"access_token"`
	Type         string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// FetchSpotifyToken retrieves a Spotify API access token in exchange for
// the authorization code that comes back upon successfully logging
// in and granting Termify access to one's Spotify information. This
// token is then used to make calls to the Spotify API.
func FetchSpotifyToken(code string) AccessToken {
	client := &http.Client{}
	accessToken := AccessToken{}
	req, err := http.NewRequest("POST", tokenURL, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addTokenQueryParams(req, code)
	addTokenHeaders(req)
	resp, _ := client.Do(req)
	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &accessToken)

	return accessToken
}

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
	q.Add("state", util.GenerateRandomString(32))
	q.Add("scope", scopes)
	req.URL.RawQuery = q.Encode()
}

func addTokenQueryParams(req *http.Request, code string) {
	q := req.URL.Query()
	q.Add("grant_type", tokenGrantType)
	q.Add("code", code)
	q.Add("redirect_uri", tokenRedirectURI)
	req.URL.RawQuery = q.Encode()
}

func addTokenHeaders(req *http.Request) {
	clientID, clientSecret := GetClientIDAndSecret()
	req.Header.Add("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
}
