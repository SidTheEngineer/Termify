package main

import (
	"Termify/app"
	"Termify/ui"
	"fmt"
)

var (
	appConfig = app.Config{}
	uiConfig  = ui.Config{}
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

func main() {
	app.InitUI(&appConfig, &uiConfig)
	fmt.Printf("%+v\n", uiConfig.CurrentView())
}
