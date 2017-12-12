package ui

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SidTheEngineer/Termify/auth"

	tui "github.com/gizak/termui"
)

const (
	playback    = "playback"
	welcomeText = "Welcome to Termify!\n\nL - Start Spotify authorization\nQ - Exit"
)

// Config will hold all of the meta data about the current state of
// Termify, such as the currentView, history stack, etc.
type Config struct {
	currentView View
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

var uiConfig Config

// CreateAPIRequest returns an http request pointer for the user selected
// choice object that is passed in.
func (c Choice) CreateAPIRequest(t auth.AccessToken) *http.Request {
	req, err := http.NewRequest(c.APIMethod, c.APIRoute, nil)

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

// CurrentView returns the View that Termify is currently displaying.
func (v *Config) CurrentView() View {
	return v.currentView
}

// Render updates the current view of Termify.
func (v *Config) Render(newView View) {
	resetRows()
	switch newView.Name {
	case playback:
		mountRow(playbackComponent())
	default:
		mountRow(welcomeComponent())
	}
	v.currentView = newView
}

// ResetRows resets the current ui rows that are being displayed
func resetRows() {
	tui.Body.Rows = tui.Body.Rows[:0]
}

func mountRow(component tui.GridBufferer) {
	tui.Body.AddRows(tui.NewCol(12, 0, component))
	tui.Body.Align()
	tui.Render(tui.Body)
}
