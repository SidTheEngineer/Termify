package ui

import (
	"Termify/auth"
	"fmt"
	"net/http"
	"os"

	tui "github.com/gizak/termui"
)

const (
	browse      = "browse"
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

// SetCurrentView updates the current view of Termify. It does not necessarily
// mean that the view is being disiplayed, however.
func (v *Config) SetCurrentView(newView View) {
	v.currentView = newView
}

// CreateWelcomeParagraph returns the initial welcome Par upon app start.
func CreateWelcomeParagraph() *tui.Par {
	welcomePar := tui.NewPar(welcomeText)
	welcomePar.Height = 10
	welcomePar.Border = false
	welcomePar.TextFgColor = tui.ColorGreen
	welcomePar.PaddingLeft = 2
	welcomePar.PaddingTop = 2

	return welcomePar
}

// CreateInitialChoiceList returns the initial List of choices to make
// upon app start.
func CreateInitialChoiceList() *tui.List {
	choiceList := tui.NewList()
	choiceList.Border = true
	choiceList.BorderFg = tui.ColorGreen
	choiceList.Height = 50
	choiceList.Items = []string{
		PlayChoice().Name,
		PauseChoice().Name,
	}

	return choiceList
}
