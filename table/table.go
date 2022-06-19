package table

import (
	"github.com/openmarketplaceengine/go-osrm/types"
)

// Request represents a request to the Table method.
// The Table method computes the duration of the fastest route between all pairs
// of supplied coordinates. Returns the durations or distances or both between
// the coordinate pairs. Note that the distances are not the shortest distance
// between two coordinates, but rather the distances of the fastest routes.
// Duration is in seconds and distance is in meters.
// See https://github.com/Project-OSRM/osrm-backend/blob/master/docs/http.md#table-service
type Request struct {
	Profile               string
	Coordinates           types.Coordinates
	Sources, Destinations []int
}

// Response represents a response from the table method
type Response struct {
	types.ResponseStatus
	Durations [][]float32 `json:"durations"`
}

func (r Request) Request() *types.Request {
	opts := types.Options{}
	if len(r.Sources) > 0 {
		opts.AddInt("sources", r.Sources...)
	}
	if len(r.Destinations) > 0 {
		opts.AddInt("destinations", r.Destinations...)
	}

	return &types.Request{
		Profile: r.Profile,
		Coords:  r.Coordinates,
		Service: "table",
		Options: opts,
	}
}
