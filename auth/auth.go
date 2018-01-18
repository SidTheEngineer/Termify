package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SidTheEngineer/Termify/util"
	"github.com/boltdb/bolt"
	"github.com/skratchdot/open-golang/open"
)

const (
	authURL               = "https://accounts.spotify.com/authorize"
	authRedirectURI       = "http://localhost:8000/callback"
	scopes                = "user-read-playback-state user-modify-playback-state"
	tokenURL              = "https://accounts.spotify.com/api/token"
	tokenGrantType        = "authorization_code"
	tokenRedirectURI      = "http://localhost:8000/callback"
	tokenMethod           = "POST"
	refreshTokenGrantType = "refresh_token"
	accessTokenText       = "accessToken"
	tokenTypeText         = "tokenType"
	tokenScopeText        = "tokenScope"
	refreshTokenText      = "refreshToken"
	tokenExpiresInText    = "tokenExpiresIn"
	timeTokenCachedText   = "timeTokenCached"
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

// Config is a type struct that is used to hold useful information
// that can be used throughout the application and the Spotify Web API.
type Config struct {
	// Params needed to fetch a Spotify access token after login/permission grant.
	AccessCode string
	ReqState   string
	AccessErr  string

	// Access token to be returned via fetch after login/permission grant/
	AccessToken AccessToken
}

// SetTokenFetchRequirements sets the proper fields needed to fetch an
// access token from Spotify after login/permission grant.
func (c *Config) SetTokenFetchRequirements(code, state, err string) {
	c.AccessCode = code
	c.ReqState = state
	c.AccessErr = err
}

// SetAccessToken sets the SpotifyConfig access token to be used throughout
// Spoitfy Web API endpoints.
func (c *Config) SetAccessToken(token AccessToken) {
	c.AccessToken = token
}

// FetchSpotifyToken retrieves a Spotify API access token in exchange for
// the authorization code that comes back upon successfully logging
// in and granting Termify access to one's Spotify information. This
// token is then used to make calls to the Spotify API.
func FetchSpotifyToken(code string) AccessToken {
	client := &http.Client{}
	accessToken := AccessToken{}
	req, err := http.NewRequest(tokenMethod, tokenURL, nil)

	if err != nil {
		log.Fatal(err)
	}

	addTokenQueryParams(req, code)
	addTokenHeaders(req)
	resp, _ := client.Do(req)
	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &accessToken)

	return accessToken
}

// FetchSpotifyTokenByRefresh uses a Spotify refresh token that was returned
// from an initial token fetch to get a new AccessToken. Used once a token has
// expired.
func FetchSpotifyTokenByRefresh(refreshToken string) AccessToken {
	client := &http.Client{}
	accessToken := AccessToken{}
	req, err := http.NewRequest(tokenMethod, tokenURL, nil)

	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("grant_type", refreshTokenGrantType)
	q.Add("refresh_token", refreshToken)
	req.URL.RawQuery = q.Encode()
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

// TokenIsExpired uses a tokens attributes to return whether or not
// "expiresIn" amount of time has passed.
func TokenIsExpired(timeTokenCached, expiresIn string) bool {
	cacheTime, _ := strconv.Atoi(timeTokenCached)
	secsTilExpire, _ := strconv.Atoi(expiresIn)

	return int(time.Now().Unix())-cacheTime > secsTilExpire
}

// CacheToken stores appropriate spotify AccessToken information in
// the bolt database
func CacheToken(tx *bolt.Tx, token AccessToken) {
	authBucket := tx.Bucket([]byte("auth"))
	authBucket.Put([]byte(accessTokenText), []byte(token.Token))
	authBucket.Put([]byte(tokenTypeText), []byte(token.Type))
	authBucket.Put([]byte(tokenScopeText), []byte(token.Scope))
	authBucket.Put([]byte(refreshTokenText), []byte(token.RefreshToken))
	authBucket.Put([]byte(tokenExpiresInText), []byte(strconv.Itoa(token.ExpiresIn)))
	authBucket.Put([]byte(timeTokenCachedText), []byte(strconv.FormatInt(int64(time.Now().Unix()), 10)))
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

	if clientID == "" {
		log.Fatal("\n\nEnvironment variable SPOTIFY_CLIENT not set.")
	}

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

	if clientID == "" || clientSecret == "" {
		log.Fatal("\n\nEnvironment variable SPOTIFY_CLIENT or SPOTIFY_SECRET not set.")
	}

	req.Header.Add("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
}
