package api

import (
	"Termify/helpers"
	"Termify/ui"
)

// Device is a response object from respective Spotify API endpoints.
type Device struct {
	id            *string
	isActive      bool
	isRestricted  bool
	name          string
	deviceType    string
	volumePercent *float64
	responseView  ui.View
}

// ResponseView gets the view for a specified Response object
func (d Device) ResponseView() ui.View {
	return d.responseView
}

// DevicesChoice returns a Choice corresponding to the Spotify "devices" endpoint.
// https://developer.spotify.com/web-api/get-a-users-available-devices/
func DevicesChoice() ui.Choice {
	return ui.Choice{
		Name:         "Devices",
		APIRoute:     "https://api.spotify.com/v1/me/player/devices",
		APIMethod:    "GET",
		ResponseType: "devices",
	}
}

func handleDevicesEndpointResponse(jsonMap map[string]interface{}) {
	var deviceList []Device
	deviceMapList := jsonMap["devices"].([]interface{})

	for _, device := range deviceMapList {
		// These Device fields can be null
		var volumePercentVal *float64
		var idVal *string

		jsonFieldMap := device.(map[string]interface{})

		if vpAssertion, ok := jsonFieldMap["volume_percent"].(float64); ok {
			volumePercentVal = &vpAssertion
		}

		if idAssertion, ok := jsonFieldMap["id"].(string); ok {
			idVal = &idAssertion
		}

		deviceList = append(deviceList, Device{
			id:            idVal,
			isActive:      jsonFieldMap["is_active"].(bool),
			isRestricted:  jsonFieldMap["is_restricted"].(bool),
			name:          jsonFieldMap["name"].(string),
			deviceType:    jsonFieldMap["type"].(string),
			volumePercent: volumePercentVal,
			responseView:  ui.View{},
		})
	}

	// After retrieving the response data
	helpers.ClearTerm()
	deviceList[0].ResponseView().Print()
}
