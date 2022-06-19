package match

import (
	"github.com/openmarketplaceengine/go-osrm/route"
	"github.com/openmarketplaceengine/go-osrm/types"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// Request represents a request to the match method.
// See https://github.com/Project-OSRM/osrm-backend/blob/master/docs/http.md#match-service
type Request struct {
	// Mode of transportation, is determined statically by the Lua profile that
	// is used to prepare the data using osrm-extract. Typically, car, bike or
	// foot if using one of the supplied profiles.
	Profile string

	// String of format
	// {longitude},{latitude};{longitude},{latitude}[;{longitude},{latitude} ...]
	// or polyline
	// ({polyline}) or polyline6({polyline6})
	Coordinates types.Coordinates

	// Limits the search to segments with given bearing in degrees towards true
	// north in clockwise direction.
	Bearings []types.Bearing

	// Returned route steps for each route
	Steps types.Steps

	// Returns additional metadata for each coordinate along the route geometry.
	Annotations types.Annotations

	// Tidy allows the input track modification to obtain better matching
	// quality for noisy tracks.
	Tidy types.Tidy

	// Timestamps for the input locations in seconds since UNIX epoch.
	// Timestamps need to be monotonically increasing.
	Timestamps []int64

	// Standard deviation of GPS precision used for map matching. If applicable
	// use GPS accuracy.
	Radii []float64

	// Hint from previous request to derive position in street network.
	Hints []string

	// Add overview geometry either full, simplified according to highest zoom
	// level it could be display on, or not at all.
	Overview types.Overview

	// Allows the input track splitting based on huge timestamp gaps between
	// points.
	Gaps types.Gaps

	// Returned route geometry format (influences overview and per step)
	Geometries types.Geometries
}

// Response represents a response from the match method
type Response struct {
	types.ResponseStatus

	// An array of Route objects that assemble the trace
	Matchings []Matching `json:"matchings"`

	// Array of Waypoint objects representing all points of the trace in order.
	// If the trace point was ommited by map matching because it is an outlier,
	// the entry will be null. Each Waypoint object has the following additional
	// properties:
	Tracepoints []*Tracepoint `json:"tracepoints"`
}

// Matching represents an array of Route objects that assemble the trace
type Matching struct {
	route.Route
	// Confidence of the matching. float value between 0 and 1. 1 is very
	// confident that the matching is correct.
	Confidence float64            `json:"confidence"`
	Geometry   geojson.LineString `json:"geometry"`
}

func (r Request) Request() *types.Request {
	options := matcherOptions(
		route.StepsOptions(r.Steps, r.Annotations, r.Overview, r.Geometries),
		r.Tidy,
		r.Gaps,
	)
	if len(r.Timestamps) > 0 {
		options.AddInt64("timestamps", r.Timestamps...)
	}
	if len(r.Radii) > 0 {
		options.AddFloat("radiuses", r.Radii...)
	}
	if len(r.Hints) > 0 {
		options.Add("hints", r.Hints...)
	}
	if len(r.Bearings) > 0 {
		options.Set("bearings", types.Bearings(r.Bearings))
	}

	return &types.Request{
		Profile: r.Profile,
		Coords:  r.Coordinates,
		Service: "match",
		Options: options,
	}
}

// Tracepoint represents a matched point on a route
type Tracepoint struct {
	// Index of the waypoint inside the matched route.
	Index    int       `json:"waypoint_index"`
	Location orb.Point `json:"location"`

	//Index to the Route object in matchings the sub-trace was matched to.
	MatchingIndex int `json:"matchings_index"`

	// Number of probable alternative matchings for this trace point. A value of
	// zero indicate that this point was matched unambiguously. Split the trace
	// at these points for incremental map matching.
	AlternativesCount int    `json:"alternatives_count"`
	Hint              string `json:"hint"`
}

func matcherOptions(options types.Options, tidy types.Tidy, gaps types.Gaps) types.Options {
	return options.
		SetStringer("tidy", tidy).
		SetStringer("gaps", gaps)
}
