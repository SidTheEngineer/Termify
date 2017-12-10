package main

import (
	"Termify/api"
	"Termify/app"
	"Termify/auth"
	"Termify/helpers"
	"Termify/server"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const (
	invalidChoiceError    = "Invalid choice. Try again from the numbered choices above."
	stringConversionError = "An error ocurred when parsing the text that was inputted. Try again."
	startupText           = "termify"
	// TODO: Set exit code programatically (could be the last choice in the view choice list).
	exitCode = 9
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
			color.Red(fmt.Sprint(stringConversionError))
			continue
		} else if inputChoice != exitCode && inputChoice > len(appConfig.CurrentView().Choices())-1 {
			color.Red(fmt.Sprint(invalidChoiceError))
			continue
		} else if inputChoice == exitCode {
			return
		}

		selectedChoice := appConfig.CurrentView().Choices()[inputChoice]
		apiReq := selectedChoice.CreateAPIRequest(apiConfig.AccessToken)
		response := selectedChoice.SendAPIRequest(apiReq)
		bytes, _ := ioutil.ReadAll(response.Body)

		var responseObject interface{}

		jsonErr := json.Unmarshal(bytes, &responseObject)

		if jsonErr != nil {
			fmt.Println(err)
		}

		fmt.Println(responseObject)

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
