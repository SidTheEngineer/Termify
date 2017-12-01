package app

// Choice is what a user selects from the view menu, each of which have their own
// Spotify api routes that they hit upon selection.
type Choice struct {
	name      string
	apiRoute  string
	apiMethod string
}

// Devices is the Spotify devices endpoint choice, which returns a
// list of Spotify device objects.
// https://developer.spotify.com/web-api/get-a-users-available-devices/
var Devices = Choice{
	name:      "Devices",
	apiRoute:  "https://api.spotify.com/v1/me/player/devices",
	apiMethod: "GET",
}

// Categories is a Spotify categories endpoint choice, which returns a
// list of categories used to tag items in Spotify (on, for example, the Spotify
// player’s “Browse” tab)
// https://developer.spotify.com/web-api/get-list-categories/
var Categories = Choice{
	name:      "Categories",
	apiRoute:  "https://api.spotify.com/v1/browse/categories",
	apiMethod: "GET",
}

// FeaturedPlaylists is a Spotify featured playlists enpoint choice, which returns
// a list of Spotify featured playlists (shown, for example, on a Spotify player’s “Browse” tab).
// https://developer.spotify.com/web-api/get-list-featured-playlists/
var FeaturedPlaylists = Choice{
	name:      "Featured Playlists",
	apiRoute:  "https://api.spotify.com/v1/browse/featured-playlists",
	apiMethod: "GET",
}
