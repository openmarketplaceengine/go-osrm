package nearest

import (
	"github.com/openmarketplaceengine/go-osrm/types"
	"github.com/paulmach/orb"
)

// Request represents a request to the nearest method.
// See https://github.com/Project-OSRM/osrm-backend/blob/master/docs/http.md#nearest-service
type Request struct {
	// Mode of transportation, is determined statically by the Lua profile that
	// is used to prepare the data using osrm-extract. Typically, car, bike or
	// foot if using one of the supplied profiles.
	Profile string

	// String of format
	// {longitude},{latitude};{longitude},{latitude}[;{longitude},{latitude} ...]
	// or polyline
	// ({polyline}) or polyline6({polyline6})
	// Length should be 1.
	Coordinates types.Coordinates

	// Limits the search to segments with given bearing in degrees towards true
	// north in clockwise direction.
	Bearings []types.Bearing

	// Number of nearest matches (segments) that should be returned.
	// Defaults to 1. Should be 1 or greater.
	Number int
}

// Response represents a response from the nearest method
type Response struct {
	types.ResponseStatus
	Waypoints []Waypoint `json:"waypoints"`
}

// Waypoint represents a nearest point on a nearest query
type Waypoint struct {
	Location orb.Point `json:"location"`
	Distance float64   `json:"distance"`
	Name     string    `json:"name"`
	Hint     string    `json:"hint"`
	// Array of OpenStreetMap node ids.
	Nodes []uint64 `json:"nodes"`
}

func (r Request) Request() *types.Request {
	opts := types.Options{}
	if r.Number > 0 {
		opts.AddInt("number", r.Number)
	}

	if len(r.Bearings) > 0 {
		opts.Set("bearings", types.Bearings(r.Bearings))
	}

	return &types.Request{
		Profile: r.Profile,
		Service: "nearest",
		Coords:  r.Coordinates,
		Options: opts,
	}
}
