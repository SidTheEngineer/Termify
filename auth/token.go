package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
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
