package table

import (
	"github.com/openmarketplaceengine/go-osrm/types"
)

// Request represents a request to the table method
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
