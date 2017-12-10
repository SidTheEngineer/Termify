package api

// Device is a response object from respective Spotify API endpoints.
type Device struct {
	ID            string   `json:"id"`
	IsActive      bool     `json:"is_active"`
	IsRestricted  bool     `json:"is_restricted"`
	Name          string   `json:"name"`
	VolumePercent int      `json:"volume_percent"`
	Devices       []Device `json:"devices"`
}
