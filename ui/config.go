package ui

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SidTheEngineer/Termify/auth"

	tui "github.com/gizak/termui"
)

const (
	// ExitText describes global exit text for the app
	ExitText = "[ Q ] - Exit"

	// NewLine can be used in termui lists and other components to be an "empty" text row
	NewLine = "\n"

	playback    = "playback"
	welcomeText = "Welcome to Termify!\n\nL - Start Spotify authorization\n" + ExitText
)

// Config will hold all of the meta data about the current state of
// Termify, such as the currentView, history stack, etc.
type Config struct {
	currentView View
	AccessToken auth.AccessToken
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

// CreateAPIRequest returns an http request pointer for the user selected
// choice object that is passed in.
func (c Choice) CreateAPIRequest(t auth.AccessToken) *http.Request {
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
	tui.Close()
	tui.Init()
}

func mountRow(component tui.GridBufferer) {
	test := tui.NewPar("this is a test component")
	test.Height = 10
	test.Border = true
	tui.Body.AddRows(tui.NewRow(
		tui.NewCol(2, 0, component),
		tui.NewCol(10, 0, test),
	))
	tui.Body.Align()
	tui.Render(tui.Body)
}
