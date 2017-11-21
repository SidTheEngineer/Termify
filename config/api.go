package config

import "Termify/auth"

// APIInformation is a type struct that is used to hold useful information
// that can be used throughout the application and the Spotify Web API.
type APIInformation struct {
	// Params needed to fetch a Spotify access token after login/permission grant.
	AccessCode string
	ReqState   string
	AccessErr  string

	// Access token to be returned via fetch after login/permission grant/
	AccessToken *auth.AccessToken
}

// SetTokenFetchRequirements sets the proper fields needed to fetch an
// access token from Spotify after login/permission grant.
func (sc *APIInformation) SetTokenFetchRequirements(code, state, err string) {
	sc.AccessCode = code
	sc.ReqState = state
	sc.AccessErr = err
}

// SetAccessToken sets the SpotifyConfig access token to be used throughout
// Spoitfy Web API endpoints.
func (sc *APIInformation) SetAccessToken(t *auth.AccessToken) {
	sc.AccessToken = t
}
