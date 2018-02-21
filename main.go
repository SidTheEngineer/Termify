package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SidTheEngineer/Termify/db"

	"github.com/SidTheEngineer/Termify/auth"
	tdb "github.com/SidTheEngineer/Termify/db"
	"github.com/SidTheEngineer/Termify/playbackUI"
	"github.com/SidTheEngineer/Termify/util"
	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	tui "github.com/gizak/termui"
)

const (
	port             = ":8000"
	grantAccessError = "A Spotfiy permission error occurred. Try logging in again."
	dbName           = "SpotfiyAuth.db"
)

var (
	authConfig auth.Config
	uiConfig   playbackUI.Config
)

func startServer(srv *http.Server) {
	srv.ListenAndServe()
}

func createServer(authConfig *auth.Config, uiConfig *playbackUI.Config) *http.Server {
	srv := &http.Server{Addr: port}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, srv, authConfig, uiConfig)
	})

	return srv
}

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server, authConfig *auth.Config, uiConfig *playbackUI.Config) {
	authConfig.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"),
	)

	if authConfig.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		token := auth.FetchSpotifyToken(authConfig.AccessCode)
		authConfig.SetAccessToken(token)
		uiConfig.SetAccessToken(token)
		// Cache the token info
		tdb.DB.Batch(func(tx *bolt.Tx) error {
			defer s.Close()
			db.CacheAccessToken(tx, token)
			return nil
		})
		util.ResetTerminal()
		playbackUI.NewPlaybackComponent(uiConfig).Render(uiConfig)
	}
}

func attachLoginHandlers() {
	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	tui.Handle("/sys/kbd/l", func(tui.Event) {
		go auth.Authorize()
		srv := createServer(&authConfig, &uiConfig)
		startServer(srv)
	})
}

func main() {
	tdb.Start()
	defer tdb.Close()

	isLoggedIn := tdb.IsLoggedIn(&uiConfig, &authConfig)

	if !isLoggedIn {
		if err := tui.Init(); err != nil {
			log.Fatal(err)
		}

		defer tui.Close()
		playbackUI.NewWelcomeComponent().Render(&uiConfig)
		attachLoginHandlers()
		tui.Loop()
	} else {
		if err := tui.Init(); err != nil {
			log.Fatal(err)
		}

		defer tui.Close()
		playbackUI.NewPlaybackComponent(&uiConfig).Render(&uiConfig)
		tui.Loop()
	}
}
