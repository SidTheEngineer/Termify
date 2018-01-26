package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SidTheEngineer/Termify/auth"
	"github.com/SidTheEngineer/Termify/playbackUI"
	"github.com/SidTheEngineer/Termify/util"
	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	tui "github.com/gizak/termui"
)

const (
	port                = ":8000"
	grantAccessError    = "A Spotfiy permission error occurred. Try logging in again."
	accessTokenText     = "accessToken"
	tokenTypeText       = "tokenType"
	tokenScopeText      = "tokenScope"
	refreshTokenText    = "refreshToken"
	tokenExpiresInText  = "tokenExpiresIn"
	timeTokenCachedText = "timeTokenCached"
	dbName              = "SpotfiyAuth.db"
)

var (
	authConfig auth.Config
	uiConfig   playbackUI.Config
	db         *bolt.DB
)

func startDB() {
	newDb, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 5 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	db = newDb
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("auth"))
		return nil
	})
}

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
		db.Batch(func(tx *bolt.Tx) error {
			defer s.Close()
			auth.CacheToken(tx, token)
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

func isLoggedIn() bool {
	loggedIn := true
	db.Batch(func(tx *bolt.Tx) error {
		authBucket := tx.Bucket([]byte("auth"))
		accessToken := authBucket.Get([]byte(accessTokenText))
		tokenType := authBucket.Get([]byte(tokenTypeText))
		tokenScope := authBucket.Get([]byte(tokenScopeText))
		refreshToken := authBucket.Get([]byte(refreshTokenText))
		expiresIn := authBucket.Get([]byte(tokenExpiresInText))
		timeTokenCached := authBucket.Get([]byte(timeTokenCachedText))

		allFieldsCached := !util.IsNil(
			accessToken,
			tokenType,
			refreshToken,
			expiresIn,
			timeTokenCached,
		) && !util.IsEmpty(
			accessToken,
			tokenType,
			refreshToken,
			expiresIn,
			timeTokenCached,
		)

		if allFieldsCached {
			if auth.TokenIsExpired(string(timeTokenCached), string(expiresIn)) {
				token := auth.FetchSpotifyTokenByRefresh(string(refreshToken))
				auth.CacheToken(tx, token)
				uiConfig.SetAccessToken(token)
				authConfig.SetAccessToken(token)
			} else {
				expireTime, _ := strconv.Atoi(string(expiresIn))

				uiConfig.SetAccessToken(auth.AccessToken{
					Token:        string(accessToken),
					Type:         string(tokenType),
					Scope:        string(tokenScope),
					RefreshToken: string(refreshToken),
					ExpiresIn:    expireTime,
				})
				authConfig.SetAccessToken(auth.AccessToken{
					Token:        string(accessToken),
					Type:         string(tokenType),
					Scope:        string(tokenScope),
					RefreshToken: string(refreshToken),
					ExpiresIn:    expireTime,
				})
			}
			// TODO: Return error that might be generated
			return nil
		}
		loggedIn = false
		return nil
	})
	return loggedIn
}

func main() {
	startDB()
	defer db.Close()

	isLoggedIn := isLoggedIn()

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
