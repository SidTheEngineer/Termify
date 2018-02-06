package playbackUI

import (
	"strconv"
	"time"

	tui "github.com/gizak/termui"
)

const (
	controlsBorderLabel      = "Controls"
	controlsWidth            = 5
	controlsHeight           = 10
	playChoiceNameText       = "[ 1 ] - Play"
	pauseChoiceNameText      = "[ 2 ] - Pause"
	previousChoiceNameText   = "[ 3 ] - Previous"
	nextChoiceNameText       = "[ 4 ] - Next"
	volumeDownChoiceNameText = "[ 5 ] - Vol. Down"
	volumeUpChoiceNameText   = "[ 6 ] - Vol. Up"
	playSuccessText          = "Play req sent"
	pauseSuccessText         = "Pause req sent"
	playErrorText            = "ERROR: Failed to send play req"
	pauseErrorText           = "ERROR: Failed to send pause req"
)

func createControls(uiConfig *Config) *tui.List {
	controls := tui.NewList()
	controls.Border = true
	controls.BorderFg = themeBorderFg
	controls.BorderLabel = controlsBorderLabel
	controls.Height = controlsHeight
	controls.ItemFgColor = themeTextFgColor
	controls.Items = []string{
		exitText,
		newLine,
		playChoice().Name,
		pauseChoice().Name,
		skipChoice().Name,
		backChoice().Name,
		volumeDownChoice(uiConfig).Name,
		volumeUpChoice(uiConfig).Name,
	}

	tui.Handle(quitKey, func(tui.Event) {
		tui.StopLoop()
	})

	attachControlsHandlers(uiConfig)

	return controls
}

func attachControlsHandlers(uiConfig *Config) {
	playbackChoices := NewPlaybackComponent(uiConfig).view.Choices

	// Unfortunately, these have to be hardcoded. Handle() breaks when trying to
	// attach in a loop.
	tui.Handle(playKey, func(e tui.Event) {
		req := playbackChoices[0].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[0].SendAPIRequest(req)
		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(updateUIWaitTime * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
			updateMessageBox(uiConfig, playSuccessText)
		} else {
			// TODO: Handle/display errors with message box
			updateMessageBox(uiConfig, playErrorText)
		}
	})

	tui.Handle(pauseKey, func(e tui.Event) {
		req := playbackChoices[1].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[1].SendAPIRequest(req)
		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(updateUIWaitTime * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
			updateMessageBox(uiConfig, pauseSuccessText)
		} else {
			// TODO: Handle/display errors with the message box
			updateMessageBox(uiConfig, pauseErrorText)
		}
	})

	tui.Handle(prevKey, func(e tui.Event) {
		req := playbackChoices[2].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[2].SendAPIRequest(req)
		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(updateUIWaitTime * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
		}
	})

	tui.Handle(nextKey, func(e tui.Event) {
		req := playbackChoices[3].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[3].SendAPIRequest(req)
		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(updateUIWaitTime * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
		}
	})

	tui.Handle(volDownKey, func(e tui.Event) {
		req := volumeDownChoice(uiConfig).CreateAPIRequest(uiConfig.AccessToken)
		res := volumeDownChoice(uiConfig).SendAPIRequest(req)
		if res.StatusCode == 204 {
			time.Sleep(updateUIWaitTime * time.Millisecond)

			// TODO: Check for < 0 or > 100
			updateVolumeGauge(uiConfig, -10)
			uiConfig.currentDevice.Volume = uiConfig.currentDevice.Volume - 10
		}
	})

	tui.Handle(volUpKey, func(e tui.Event) {
		req := volumeUpChoice(uiConfig).CreateAPIRequest(uiConfig.AccessToken)
		res := volumeUpChoice(uiConfig).SendAPIRequest(req)
		if res.StatusCode == 204 {
			time.Sleep(updateUIWaitTime * time.Millisecond)

			// TODO: Check for < 0 or > 100
			updateVolumeGauge(uiConfig, 10)
			uiConfig.currentDevice.Volume = uiConfig.currentDevice.Volume + 10
		}
	})
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/start-a-users-playback/
func playChoice() Choice {
	return Choice{
		Name:         playChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/play",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/pause-a-users-playback/
func pauseChoice() Choice {
	return Choice{
		Name:         pauseChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/pause",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-next-track/
func skipChoice() Choice {
	return Choice{
		Name:         previousChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/previous",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-previous-track/
func backChoice() Choice {
	return Choice{
		Name:         nextChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/next",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/set-volume-for-users-playback/
func volumeDownChoice(uiConfig *Config) Choice {
	apiRoute := "https://api.spotify.com/v1/me/player/volume?volume_percent="
	newVolume := strconv.Itoa(int(uiConfig.currentDevice.Volume - 10))
	deviceIDParam := "&device_id=" + uiConfig.currentDevice.ID
	return Choice{
		Name:         volumeDownChoiceNameText,
		APIRoute:     apiRoute + newVolume + deviceIDParam,
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

func volumeUpChoice(uiConfig *Config) Choice {
	apiRoute := "https://api.spotify.com/v1/me/player/volume?volume_percent="
	newVolume := strconv.Itoa(int(uiConfig.currentDevice.Volume + 10))
	deviceIDParam := "&device_id=" + uiConfig.currentDevice.ID
	return Choice{
		Name:         volumeUpChoiceNameText,
		APIRoute:     apiRoute + newVolume + deviceIDParam,
		APIMethod:    "PUT",
		ResponseType: "",
	}
}
