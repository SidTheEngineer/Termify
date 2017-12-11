package app

import (
	"Termify/auth"
	"Termify/ui"
	"fmt"
	"net/http"
	"os"

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
func InitUI(appConfig Config, uiConfig ui.Config) {
	if err := tui.Init(); err != nil {
		panic(err)
	}
	defer tui.Close()

	tui.Body.AddRows(tui.NewRow(tui.NewCol(12, 0, ui.CreateWelcomeParagraph())))

	tui.Body.Align()
	tui.Render(tui.Body)

	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	tui.Handle("/sys/kbd/l", func(tui.Event) {
		go auth.Authorize()
		srv := createServer(&appConfig, &uiConfig)
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

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server, apiConfig *Config, uiConfig *ui.Config) {
	defer s.Close()

	apiConfig.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"))

	if apiConfig.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		apiConfig.SetAccessToken(auth.FetchSpotifyToken(apiConfig.AccessCode))

		choiceList := tui.NewList()
		choiceList.Height = 30
		choiceList.BorderLabel = "Termify"
		choiceList.BorderFg = tui.ColorCyan
		choiceList.Items = []string{
			ui.PlayChoice().Name,
			ui.PauseChoice().Name,
		}

		tui.Body.Rows = tui.Body.Rows[:0]
		tui.Body.AddRows(tui.NewRow(tui.NewCol(12, 0, ui.CreateInitialChoiceList())))
		tui.Body.Align()
		tui.Render(tui.Body)
		view := ui.View{
			Name: "playback",
			Choices: []ui.Choice{
				ui.PlayChoice(),
				ui.PauseChoice(),
			},
		}
		uiConfig.SetCurrentView(view)
	}
}
