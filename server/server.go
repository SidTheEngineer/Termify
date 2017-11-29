package server

import (
	"Termify/api"
	"Termify/app"
	"Termify/auth"
	"Termify/helpers"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/fatih/color"
)

const (
	port               = ":8000"
	grantAccessError   = "A Spotfiy permission error occurred. Try logging in again. "
	parseTemplateError = "An error occurred when trying to parse a template"
)

var (
	apiConfig = api.Config{}
	appConfig = app.Config{}
)

// Create creates a server instance on some supplied global port,
// attaches necessary handlers, and returns the server to be used
// elsewhere.
func Create() *http.Server {
	srv := &http.Server{Addr: port}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, srv)
	})
	return srv
}

// Start starts a provided http Server instance.
func Start(srv *http.Server) {
	srv.ListenAndServe()
}

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server) {
	apiConfig.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"))

	if apiConfig.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		app.Init()
		apiConfig.SetAccessToken(auth.FetchSpotifyToken(apiConfig.AccessCode))

		// The first view that needs to display is the browse view.
		view := app.NewView("browse")

		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)

		helpers.ClearTerm()
		view.Print()
		defer s.Close()
	}
}
