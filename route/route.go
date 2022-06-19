package route

import (
	"fmt"
	"github.com/openmarketplaceengine/go-osrm/types"
	"github.com/paulmach/orb/geojson"
	"strconv"

	"github.com/paulmach/orb"
)

// Request represents a request to the Route method.
// The Route service finds the fastest route between coordinates in the supplied
// order.
type Request struct {
	Profile     string
	Coordinates types.Coordinates
	Bearings    []types.Bearing

	// Returned route steps for each route leg
	Steps types.Steps

	// Returns additional metadata for each coordinate along the route geometry.
	Annotations types.Annotations

	// Add overview geometry either full, simplified according to highest zoom
	// level it could be display on, or not at all.
	Overview types.Overview

	// Returned route geometry format (influences overview and per step)
	Geometries types.Geometries

	// Forces the route to keep going straight at waypoints constraining uturns
	// there even if it would be faster. Default value depends on the profile.
	ContinueStraight types.ContinueStraight
	Waypoints        []int
}

// Response represents a response from the route method
type Response struct {
	types.ResponseStatus
	Routes    []Route    `json:"routes"`
	Waypoints []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	Name     string    `json:"name"`
	Location orb.Point `json:"location"`
	Distance float32   `json:"distance"`
	Hint     string    `json:"hint"`
}

// Route represents a route through (potentially multiple) points.
type Route struct {
	Distance   float32            `json:"distance"`
	Duration   float32            `json:"duration"`
	WeightName string             `json:"weight_name"`
	Weight     float32            `json:"weight"`
	Geometry   geojson.LineString `json:"geometry"`
	Legs       []Leg              `json:"legs"`
}

// Leg represents a route between two waypoints.
type Leg struct {
	Annotation Annotation `json:"annotation"`
	Distance   float32    `json:"distance"`
	Duration   float32    `json:"duration"`
	Summary    string     `json:"summary"`
	Weight     float32    `json:"weight"`
	Steps      []Step     `json:"steps"`
}

// Annotation contains additional metadata for each coordinate along the route geometry
type Annotation struct {
	Duration []float32 `json:"duration,omitempty"`
	Distance []float32 `json:"distance,omitempty"`
	Nodes    []uint64  `json:"nodes,omitempty"`
}

// Step represents a route geometry
type Step struct {
	Distance      float32            `json:"distance"`
	Duration      float32            `json:"duration"`
	Geometry      geojson.LineString `json:"geometry"`
	Name          string             `json:"name"`
	Mode          string             `json:"mode"`
	DrivingSide   string             `json:"driving_side"`
	Weight        float32            `json:"weight"`
	Maneuver      StepManeuver       `json:"maneuver"`
	Intersections []Intersection     `json:"intersections,omitempty"`
}

type Intersection struct {
	Location orb.Point `json:"location"`
	Bearings []uint16  `json:"bearings"`
	Entry    []bool    `json:"entry"`
	In       *uint32   `json:"in,omitempty"`
	Out      *uint32   `json:"out,omitempty"`
	Lanes    []Lane    `json:"lanes,omitempty"`
}

type Lane struct {
	Indications []string `json:"indications"`
	Valid       bool     `json:"valid"`
}

// StepManeuver contains information about maneuver in step
type StepManeuver struct {
	Location      orb.Point `json:"location"`
	BearingBefore float32   `json:"bearing_before"`
	BearingAfter  float32   `json:"bearing_after"`
	Type          string    `json:"type"`
	Modifier      string    `json:"modifier,omitempty"`
	Exit          *uint32   `json:"exit,omitempty"`
}

func (r Request) Request() *types.Request {
	opts := StepsOptions(r.Steps, r.Annotations, r.Overview, r.Geometries).
		SetStringer("continue_straight", r.ContinueStraight)

	if len(r.Waypoints) > 0 {
		waypoints := ""
		for i, w := range r.Waypoints {
			if i > 0 {
				waypoints += ";"
			}
			waypoints += strconv.Itoa(w)
		}
		opts.Set("waypoints", waypoints)
	}

	if len(r.Bearings) > 0 {
		opts.Set("bearings", types.Bearings(r.Bearings))
	}

	return &types.Request{
		Profile: r.Profile,
		Coords:  r.Coordinates,
		Service: "route",
		Options: opts,
	}
}

func StepsOptions(steps types.Steps, annotations types.Annotations, overview types.Overview, geometries types.Geometries) types.Options {
	return types.Options{}.
		SetStringer("steps", steps).
		SetStringer("annotations", annotations).
		SetStringer("geometries", valueOrDefault(geometries, types.GeometriesPolyline6)).
		SetStringer("overview", overview)
}

func valueOrDefault(value, def fmt.Stringer) fmt.Stringer {
	if value.String() == "" {
		return def
	}
	return value
}
