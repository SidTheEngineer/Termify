package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SidTheEngineer/Termify/auth"
	"github.com/SidTheEngineer/Termify/ui"
	"github.com/fatih/color"
	tui "github.com/gizak/termui"
)

const (
	port             = ":8000"
	grantAccessError = "A Spotfiy permission error occurred. Try logging in again."
)

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

// InitUI initializes the UI and attaches the appropriate handlers.
func InitUI(appConfig *Config, uiConfig *ui.Config) {
	if err := tui.Init(); err != nil {
		panic(err)
	}
	defer tui.Close()

	uiConfig.Render(ui.View{
		Name: "welcome",
	})

	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	tui.Handle("/sys/kbd/l", func(tui.Event) {
		go auth.Authorize()
		srv := createServer(appConfig, uiConfig)
		startServer(srv)
	})

	tui.Loop()
}

func createServer(appConfig *Config, uiConfig *ui.Config) *http.Server {
	srv := &http.Server{Addr: port}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, srv, appConfig, uiConfig)
	})

	return srv
}

func startServer(srv *http.Server) {
	srv.ListenAndServe()
}

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server, appConfig *Config, uiConfig *ui.Config) {
	defer s.Close()

	appConfig.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"))

	if appConfig.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		appConfig.SetAccessToken(auth.FetchSpotifyToken(appConfig.AccessCode))
		uiConfig.Render(ui.NewPlaybackView())
	}
}
