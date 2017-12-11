package main

import (
	"Termify/api"
	"Termify/app"
	"Termify/auth"
	"Termify/ui"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	tui "github.com/gizak/termui"
)

const (
	invalidChoiceError    = "Invalid choice. Try again from the numbered choices above."
	stringConversionError = "An error ocurred when parsing the text that was inputted. Try again."
	startupText           = "termify"
	port                  = ":8000"
	grantAccessError      = "A Spotfiy permission error occurred. Try logging in again. "
	parseTemplateError    = "An error occurred when trying to parse a template"
	welcomeText           = "Welcome to Termify!\n\nL - Start Spotify authorization\nQ - Exit"
)

var (
	apiConfig = api.Config{}
	appConfig = app.Config{}
)

// func startUILoop() {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		text := scanner.Text()
// 		inputChoice, err := strconv.Atoi(text)

// 		if err != nil {
// 			color.Red(fmt.Sprint(stringConversionError))
// 			continue
// 		} else if inputChoice != exitCode && inputChoice > len(appConfig.CurrentView().Choices)-1 {
// 			color.Red(fmt.Sprint(invalidChoiceError))
// 			continue
// 		} else if inputChoice == exitCode {
// 			return
// 		}

// 		selectedChoice := appConfig.CurrentView().Choices[inputChoice]
// 		apiReq := selectedChoice.CreateAPIRequest(apiConfig.AccessToken)
// 		response := selectedChoice.SendAPIRequest(apiReq)
// 		bytes, _ := ioutil.ReadAll(response.Body)

// 		var responseObject interface{}
// 		jsonErr := json.Unmarshal(bytes, &responseObject)

// 		if jsonErr != nil {
// 			fmt.Println(err)
// 		} else if responseObject == nil {
// 			continue
// 		} else {
// 			jsonMap := responseObject.(map[string]interface{})

// 			// If we have no responseType, then the request went to an endpoint
// 			// that did not return anything (it was a command or something of the like)
// 			if selectedChoice.ResponseType != "" {
// 				api.HandleJSONResponse(jsonMap, selectedChoice.ResponseType)
// 			}
// 		}
// 	}
// }

func createServer(apiConfig *api.Config, appConfig *app.Config) *http.Server {
	srv := &http.Server{Addr: port}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, srv, apiConfig, appConfig)
	})
	return srv
}

// Start starts a provided http Server instance.
func startServer(srv *http.Server) {
	srv.ListenAndServe()
}

func callbackHandler(w http.ResponseWriter, r *http.Request, s *http.Server, apiConfig *api.Config, appConfig *app.Config) {
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
			api.PlayChoice().Name,
			api.PauseChoice().Name,
		}

		tui.Body.Rows = tui.Body.Rows[:0]
		tui.Body.AddRows(tui.NewRow(tui.NewCol(12, 0, createInitialChoiceList())))
		tui.Body.Align()
		tui.Render(tui.Body)
		view := ui.View{
			Name: "playback",
			Choices: []ui.Choice{
				api.PlayChoice(),
				api.PauseChoice(),
			},
		}
		appConfig.SetCurrentView(view)
	}
}

func initUI() {
	if err := tui.Init(); err != nil {
		panic(err)
	}
	defer tui.Close()

	tui.Body.AddRows(tui.NewRow(tui.NewCol(12, 0, createWelcomeParagraph())))

	tui.Body.Align()
	tui.Render(tui.Body)

	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	tui.Handle("/sys/kbd/l", func(tui.Event) {
		go auth.Authorize()
		srv := createServer(&apiConfig, &appConfig)
		startServer(srv)
	})

	tui.Loop()
}

func createWelcomeParagraph() *tui.Par {
	welcomePar := tui.NewPar(welcomeText)
	welcomePar.Height = 10
	welcomePar.Border = false
	welcomePar.TextFgColor = tui.ColorGreen
	welcomePar.PaddingLeft = 2
	welcomePar.PaddingTop = 2

	return welcomePar
}

func createInitialChoiceList() *tui.List {
	choiceList := tui.NewList()
	choiceList.Border = true
	choiceList.BorderFg = tui.ColorGreen
	choiceList.Height = 50
	choiceList.Items = []string{
		api.PlayChoice().Name,
		api.PauseChoice().Name,
	}

	return choiceList
}

func main() {
	initUI()
}
