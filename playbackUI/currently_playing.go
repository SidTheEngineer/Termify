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
		// Stop the previous tickers if they exist
		if uiConfig.progressTicker != nil {
			uiConfig.progressTicker.Stop()
		}
		if uiConfig.visualsTicker != nil {
			uiConfig.visualsTicker.Stop()
		}

		// Adjust this according to how fast you'd want the visuals to update. BPM / 60 = how many
		// times the ticker needs to tick per second (multiplied by 1000 for milliseconds).
		visualsTickTime := time.Duration(int(1000 / (uiConfig.currentTrack.BPM / 60)))

		playingState = playingText
		uiConfig.progressTicker = time.NewTicker(time.Millisecond * 1000)
		uiConfig.visualsTicker = time.NewTicker(time.Duration(visualsTickTime * time.Microsecond * 1000))
		uiConfig.timeElapsedFromTickerStart = 0

		go startTrackProgressTicker(uiConfig, trackInfo, deviceInfo)
		go startVisualsTicker(uiConfig)
	} else {
		playingState = pausedText
		if uiConfig.progressTicker != nil {
			uiConfig.progressTicker.Stop()
		}

		if uiConfig.visualsTicker != nil {
			uiConfig.visualsTicker.Stop()
		}
	}

	currentlyPlayingUI := tui.NewList()
	currentlyPlayingUI.Border = true
	currentlyPlayingUI.BorderLabel = currentlyPlayingBorderLabel
	currentlyPlayingUI.BorderFg = tui.ColorMagenta
	currentlyPlayingUI.Height = currentlyPlayingHeight
	currentlyPlayingUI.Items = []string{
		newLine,
		deviceInfo.DeviceType + " - " + deviceInfo.Name,
		newLine + newLine,
		trackInfo.Name,
		trackInfo.Artists,
		newLine + newLine,
		playingState,
	}
	currentlyPlayingUI.ItemFgColor = themeTextFgColor

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
	currentTrack := getTrackInformationFromJSON(uiConfig, currentContext)
	deviceInfo := getDeviceInformationFromJSON(uiConfig, currentContext)
	newCurrentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, currentTrack, deviceInfo)

	tui.Body.Rows[0].Cols[1] = tui.NewCol(currentlyPlayingWidth, 0, newCurrentlyPlayingUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}

func startTrackProgressTicker(uiConfig *Config, trackInfo Track, deviceInfo Device) {
	for _ = range uiConfig.progressTicker.C {
		uiConfig.timeElapsedFromTickerStart += 1000

		if uiConfig.timeElapsedFromTickerStart+int(deviceInfo.ProgressMs) > int(trackInfo.DurationMs) {
			uiConfig.progressTicker.Stop()
			updateCurrentlyPlayingUI(uiConfig)
		}
		progressInSeconds := (uiConfig.timeElapsedFromTickerStart + int(deviceInfo.ProgressMs)) / 1000
		updateTrackProgressGauge(uiConfig, progressInSeconds)
	}
}

func startVisualsTicker(uiConfig *Config) {
	for _ = range uiConfig.visualsTicker.C {
		updatePlayingAnimationUI(uiConfig)
	}
}
