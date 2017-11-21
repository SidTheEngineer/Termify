package main

import (
	"Termify/auth"
	"Termify/config"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	port                  = ":8000"
	startupText           = "termify"
	permissionGrantedText = "\nYou have successfully logged in.\n\n"
	parseTemplateError    = "An error occurred when trying to parse a template"
	grantAccessError      = "A Spotfiy permission error occurred. Try logging in again. "
)

var apiInformation = config.APIInformation{}

func startup(text string) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	color.Green(banner.PrintS(text))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	apiInformation.SetTokenFetchRequirements(
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
		r.URL.Query().Get("error"))

	if apiInformation.AccessErr != "" {
		color.Red(fmt.Sprint(grantAccessError))
		os.Exit(1)
	} else {
		color.HiBlue(fmt.Sprint(permissionGrantedText))
		apiInformation.SetAccessToken(auth.FetchSpotifyToken(apiInformation.AccessCode))

		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/callback", callbackHandler)

	startup(startupText)
	auth.Authorize()
	defer http.ListenAndServe(port, nil)

}
