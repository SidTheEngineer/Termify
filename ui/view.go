package ui

import (
	"Termify/auth"
	"fmt"
	"net/http"
	"os"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	browse      = "browse"
	welcomeText = "Welcome to Termify!\n\nL - Start Spotify authorization\nQ - Exit"
)

// View is a struct that contains a view's information and behaviors, such
// as the type of view it is, and it being displayed in the console.
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

// Print prints a View's informatio to the console.
func (v View) Print() {
	color.Cyan(banner.PrintS(v.Name))
	fmt.Printf("\n\n")
	for i, choice := range v.Choices {
		color.HiGreen(fmt.Sprintf("%d. %s\n", i, choice.Name))
	}
	fmt.Printf("\n\n")
}

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
