package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SidTheEngineer/Termify/auth"
	"github.com/SidTheEngineer/Termify/ui"
	"github.com/boltdb/bolt"
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
	db         *bolt.DB
)

func startDB() {
	newDb, err := bolt.Open("SpotfiyAuth.db", 0600, &bolt.Options{Timeout: 5 * time.Second})

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

func createServer(authConfig *auth.Config, uiConfig *ui.Config) *http.Server {
	srv := &http.Server{Addr: port}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, srv, authConfig, uiConfig)
	})

	return srv
}

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server, authConfig *auth.Config, uiConfig *ui.Config) {
	authConfig.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"),
	)

	if authConfig.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		fmt.Println("doing the token stuff")
		token := auth.FetchSpotifyToken(authConfig.AccessCode)
		authConfig.SetAccessToken(token)
		uiConfig.SetAccessToken(token)
		db.Update(func(tx *bolt.Tx) error {
			defer s.Close()
			fmt.Println("doing the db stuff")
			authBucket := tx.Bucket([]byte("auth"))
			fmt.Println("setting the access token")
			err := authBucket.Put([]byte("accessToken"), []byte(token.Token))

			if err != nil {
				log.Fatal(err)
			}

			return nil
		})
		// uiConfig.Render(ui.NewPlaybackView(), uiConfig)
	}
}

func main() {
	startDB()
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		authBucket := tx.Bucket([]byte("auth"))
		at := authBucket.Get([]byte("accessToken"))

		if at != nil {
			fmt.Printf("the token: %s", string(at))
			return nil
		}

		if err := tui.Init(); err != nil {
			panic(err)
		}
		defer tui.Close()

		// uiConfig.Render(ui.View{
		// 	Name: "welcome",
		// }, &uiConfig)

		tui.Handle("/sys/kbd/q", func(tui.Event) {
			tui.StopLoop()
		})

		tui.Handle("/sys/kbd/l", func(tui.Event) {
			go auth.Authorize()
			srv := createServer(&authConfig, &uiConfig)
			startServer(srv)
		})

		tui.Loop()

		return nil
	})
}
