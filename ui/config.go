package ui

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SidTheEngineer/Termify/auth"

	tui "github.com/gizak/termui"
)

const (
	// ExitText describes global exit text for the app
	ExitText = "[ Q ] - Quit"

	// NewLine can be used in termui lists and other components to be an "empty" text row
	NewLine = "\n"

	playback    = "playback"
	welcomeText = "Welcome to Termify!\n\nL - Start Spotify authorization\n" + ExitText
)

// Config will hold all of the meta data about the current state of
// Termify, such as the currentView, history stack, etc.
type Config struct {
	currentView                View
	AccessToken                auth.AccessToken
	context                    map[string]interface{}
	progressTicker             *time.Ticker
	timeElapsedFromTickerStart int
	currentTrack               Track
	currentDevice              Device
}

// View is a struct that contains a view's information and behaviors, such
// as the type of view it is, and it what choices there are for that respective type.
type View struct {
	Name    string
	Choices []Choice
}

// Choice is what a user selects from the view menu, each of which have their own
// Spotify api routes that they hit upon selection.
type Choice struct {
	Name         string
	APIRoute     string
	APIMethod    string
	ResponseType string
}

// SetAccessToken sets the SpotifyConfig access token to be used throughout
// Spoitfy Web API endpoints.
func (c *Config) SetAccessToken(token auth.AccessToken) {
	c.AccessToken = token
}

// SetCurrentlyPlayingContext updates information about the current device/track
func (c *Config) SetCurrentlyPlayingContext(ctx map[string]interface{}) {
	c.context = ctx
}

// CreateAPIRequest returns an http request pointer for the user selected
// choice object that is passed in.
func (c Choice) CreateAPIRequest(t auth.AccessToken) *http.Request {
	// TODO: Find a way to get the current token and when it was cached (which
	// is held in the db) from here. May need to take a look at changing
	// architecture a bit.
	req, err := http.NewRequest(c.APIMethod, c.APIRoute, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

// CurrentView returns the View that Termify is currently displaying.
func (c *Config) CurrentView() View {
	return c.currentView
}

// ResetTerminal resets the current ui rows that are being displayed
func ResetTerminal() {
	tui.Body.Rows = tui.Body.Rows[:0]
}
