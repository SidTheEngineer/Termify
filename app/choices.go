package app

import (
	"Termify/auth"
	"fmt"
	"net/http"
	"os"
)

// Choice is what a user selects from the view menu, each of which have their own
// Spotify api routes that they hit upon selection.
type Choice struct {
	name         string
	apiRoute     string
	apiMethod    string
	responseType string
}

// Name returns the name of the specified Choice.
func (c Choice) Name() string {
	return c.name
}

// APIRoute returns the API route that corresponds to the specified Choice.
func (c Choice) APIRoute() string {
	return c.apiRoute
}

// APIMethod returns the API method that corresponds to the specified Choice/API route.
func (c Choice) APIMethod() string {
	return c.apiMethod
}

// ResponseType returns a user selected Choice's response object type from the Spotify
// API as a string.
func (c Choice) ResponseType() string {
	return c.responseType
}

// CreateAPIRequest returns an http request pointer for the user selected
// choice object that is passed in.
func (c Choice) CreateAPIRequest(t auth.AccessToken) *http.Request {
	req, err := http.NewRequest(c.APIMethod(), c.APIRoute(), nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", "Bearer "+t.Token)
	req.Header.Add("Accept", "application/json")
	return req
}

// SendAPIRequest sends a request for an API request object created from
// a user selected Choice, and returns a pointer to an http Response.
func (c Choice) SendAPIRequest(req *http.Request) *http.Response {
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resp
}

// Devices is the Spotify devices endpoint choice, which returns a
// list of Spotify device objects.
// https://developer.spotify.com/web-api/get-a-users-available-devices/
var Devices = Choice{
	name:         "Devices",
	apiRoute:     "https://api.spotify.com/v1/me/player/devices",
	apiMethod:    "GET",
	responseType: "devices",
}

// Categories is a Spotify categories endpoint choice, which returns a
// list of categories used to tag items in Spotify (on, for example, the Spotify
// player’s “Browse” tab)
// https://developer.spotify.com/web-api/get-list-categories/
var Categories = Choice{
	name:      "Categories",
	apiRoute:  "https://api.spotify.com/v1/browse/categories",
	apiMethod: "GET",
}

// FeaturedPlaylists is a Spotify featured playlists enpoint choice, which returns
// a list of Spotify featured playlists (shown, for example, on a Spotify player’s “Browse” tab).
// https://developer.spotify.com/web-api/get-list-featured-playlists/
var FeaturedPlaylists = Choice{
	name:      "Featured Playlists",
	apiRoute:  "https://api.spotify.com/v1/browse/featured-playlists",
	apiMethod: "GET",
}
