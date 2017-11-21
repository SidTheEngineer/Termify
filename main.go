package main

import (
	"Termify/auth"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	port        = ":8000"
	startupText = "termify"
)

// AccessToken is a token response returned after a Spotify authorization
// request is made.
type AccessToken struct {
	Token        string `json:"access_token"`
	Type         string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func startup(text string) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	color.Green(banner.PrintS(text))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	color.Cyan(fmt.Sprint("\nSpotify has been granted access to your profile. Welcome to Termify.\n\n"))
	fmt.Fprintf(w, "<h1>Spotify has been granted access</h1>")
}

func main() {
	startup(startupText)
	auth.Authorize()

	http.HandleFunc("/callback", callbackHandler)
	http.ListenAndServe(port, nil)
}
