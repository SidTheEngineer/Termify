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
	"strconv"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	invalidChoiceError     = "Invalid choice. Try again from the numbered choices above."
	stringConverstionError = "An error ocurred when parsing the text that was inputted. Try again."
	startupText            = "termify"
	exitCode               = 9
)

var (
	apiConfig = api.Config{}
	appConfig = app.Config{}
)

func startup(text string) {
	helpers.ClearTerm()
	color.Green(banner.PrintS(text))
}

func startUILoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		inputChoice, err := strconv.Atoi(text)

		if err != nil {
			color.Red(fmt.Sprint(stringConverstionError))
			continue
		} else if inputChoice > len(appConfig.CurrentView().Choices())-1 {
			color.Red(fmt.Sprint(invalidChoiceError))
			continue
		} else if inputChoice == exitCode {
			return
		}

		fmt.Println(appConfig.CurrentView().Choices()[inputChoice].Name())
	}
}

func main() {
	startup(startupText)
	auth.Authorize()
	srv := server.Create(&apiConfig, &appConfig)
	server.Start(srv)

	// fmt.Printf("%+v\n", apiConfig)
	// fmt.Printf("%+v\n", appConfig)

	startUILoop()
}
