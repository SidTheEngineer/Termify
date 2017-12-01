package main

import (
	"Termify/api"
	"Termify/app"
	"Termify/auth"
	"Termify/helpers"
	"Termify/server"
	"bufio"
	"fmt"
	"os"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	startupText = "termify"
	exitCode    = "9"
)

var (
	apiConfig = api.Config{}
	appConfig = app.Config{}
)

func startup(text string) {
	helpers.ClearTerm()
	color.Green(banner.PrintS(text))
}

func main() {
	startup(startupText)
	auth.Authorize()
	srv := server.Create(&apiConfig, &appConfig)
	server.Start(srv)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%+v\n", apiConfig)
	fmt.Printf("%+v\n", appConfig)

	// When on the main menu, keep scanning input until the user
	// chooses to exit
	for scanner.Scan() {
		text := scanner.Text()
		if text == exitCode {
			return
		}
		fmt.Println(text)
	}
}
