package playbackUI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	tui "github.com/gizak/termui"
)

const (
	playingText                   = "[ Playing ]"
	pausedText                    = "[ Paused ]"
	currentlyPlayingContextURL    = "https://api.spotify.com/v1/me/player"
	currentlyPlayingContextMethod = "GET"
	currentlyPlayingBorderLabel   = "Currently Playing"
	currentlyPlayingHeight        = 10
	currentlyPlayingWidth         = 7
)

func createCurrentlyPlayingUI(uiConfig *Config, trackInfo Track, deviceInfo Device) *tui.List {
	var playingState string

	if deviceInfo.IsPlaying {
		if uiConfig.progressTicker != nil {
			uiConfig.progressTicker.Stop()
		}
		playingState = playingText
		uiConfig.progressTicker = time.NewTicker(time.Millisecond * 1000)
		uiConfig.timeElapsedFromTickerStart = 0
		go func() {
			for _ = range uiConfig.progressTicker.C {
				uiConfig.timeElapsedFromTickerStart += 1000

				if uiConfig.timeElapsedFromTickerStart+int(deviceInfo.ProgressMs) > int(trackInfo.DurationMs) {
					uiConfig.progressTicker.Stop()
					updateCurrentlyPlayingUI(uiConfig)
				}
				// TODO: Calculate progress into song based on needed variables. Use uiConfig
				// where necessary.
				progressInSeconds := (uiConfig.timeElapsedFromTickerStart + int(deviceInfo.ProgressMs)) / 1000
				updatetrackProgressTime(uiConfig, progressInSeconds)
			}
		}()
		// Skipping or going back a track always plays the track as well, so we will
		// only reach this else if a pause choice is chosen.
	} else {
		playingState = pausedText
		if uiConfig.progressTicker != nil {
			uiConfig.progressTicker.Stop()
		}
	}

	currentlyPlayingUI := tui.NewList()
	currentlyPlayingUI.Border = true
	currentlyPlayingUI.BorderLabel = currentlyPlayingBorderLabel
	currentlyPlayingUI.BorderFg = tui.ColorMagenta
	currentlyPlayingUI.Height = currentlyPlayingHeight
	currentlyPlayingUI.Items = []string{
		NewLine,
		deviceInfo.DeviceType + " - " + deviceInfo.Name,
		NewLine + NewLine,
		trackInfo.Name,
		trackInfo.Artists,
		NewLine + NewLine,
		playingState,
	}
	currentlyPlayingUI.ItemFgColor = tui.ColorYellow

	return currentlyPlayingUI
}

func getCurrentlyPlayingContext(uiConfig *Config) map[string]interface{} {
	var jsonMap map[string]interface{}
	client := http.Client{}
	req, _ := http.NewRequest(currentlyPlayingContextMethod, currentlyPlayingContextURL, nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", uiConfig.AccessToken.Token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &jsonMap)

	uiConfig.SetCurrentlyPlayingContext(jsonMap)

	return jsonMap
}

func updateCurrentlyPlayingUI(uiConfig *Config) {
	currentContext := getCurrentlyPlayingContext(uiConfig)
	currentTrack := getTrackInformationFromJSON(currentContext)
	deviceInfo := getDeviceInformationFromJSON(currentContext)
	newCurrentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, currentTrack, deviceInfo)

	// Currently Playing box is row 1, column 2
	tui.Body.Rows[0].Cols[1] = tui.NewCol(currentlyPlayingWidth, 0, newCurrentlyPlayingUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}