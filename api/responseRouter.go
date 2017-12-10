package api

const (
	devices = "devices"
)

// HandleJSONResponse is a function that properly routes responses after hitting a Spotify
// endpoint and converting JSON data to a golang map.
func HandleJSONResponse(jsonMap map[string]interface{}, responseType string) {
	switch responseType {
	case devices:
		handleDevicesEndpointResponse(jsonMap)
	}
}
