package api

import "Termify/auth"

// Config is a type struct that is used to hold useful information
// that can be used throughout the application and the Spotify Web API.
type Config struct {
	// Params needed to fetch a Spotify access token after login/permission grant.
	AccessCode string
	ReqState   string
	AccessErr  string

	// Access token to be returned via fetch after login/permission grant/
	AccessToken auth.AccessToken
}

// SetTokenFetchRequirements sets the proper fields needed to fetch an
// access token from Spotify after login/permission grant.
func (info *Config) SetTokenFetchRequirements(code, state, err string) {
	info.AccessCode = code
	info.ReqState = state
	info.AccessErr = err
}

// SetAccessToken sets the SpotifyConfig access token to be used throughout
// Spoitfy Web API endpoints.
func (info *Config) SetAccessToken(token auth.AccessToken) {
	info.AccessToken = token
}
