package api

import (
	"fmt"
)

const (
	devies = "devices"
)

// Device is a response object from respective Spotify API endpoints.
type Device struct {
	ID            string   `json:"id"`
	IsActive      bool     `json:"is_active"`
	IsRestricted  bool     `json:"is_restricted"`
	Name          string   `json:"name"`
	VolumePercent int      `json:"volume_percent"`
	Devices       []Device `json:"devices"`
}

// HandleJSONResponse is a function that properly routes responses after hitting a Spotify
// endpoint and converting JSON data to a golang map.
func HandleJSONResponse(jsonMap map[string]interface{}, responseType string) {
	switch responseType {
	case devies:
		handleDevicesEndpointResponse(jsonMap)
	}
}

func handleDevicesEndpointResponse(jsonMap map[string]interface{}) {
	fmt.Println(jsonMap)
}
