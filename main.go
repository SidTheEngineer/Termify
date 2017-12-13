package main

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

var (
	authConfig auth.Config
	uiConfig   ui.Config
)

func startServer(srv *http.Server) {
	srv.ListenAndServe()
}

func createServer(authConfig *auth.Config, uiConfig *ui.Config) *http.Server {
	srv := &http.Server{Addr: port}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, srv, authConfig, uiConfig)
	})

	return srv
}

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server, authConfig *auth.Config, uiConfig *ui.Config) {
	defer s.Close()

	authConfig.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"))

	if authConfig.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		token := auth.FetchSpotifyToken(authConfig.AccessCode)
		authConfig.SetAccessToken(token)
		uiConfig.SetAccessToken(token)
		uiConfig.Render(ui.NewPlaybackView(), uiConfig)
	}
}

func main() {
	if err := tui.Init(); err != nil {
		panic(err)
	}
	defer tui.Close()

	uiConfig.Render(ui.View{
		Name: "welcome",
	}, &uiConfig)

	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	tui.Handle("/sys/kbd/l", func(tui.Event) {
		go auth.Authorize()
		srv := createServer(&authConfig, &uiConfig)
		startServer(srv)
	})

	tui.Loop()
}
