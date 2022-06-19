package types

import (
	"fmt"
	"net/url"
	"strings"
)

// Tidy represents a tidy param for osrm5 match request
type Tidy string

// Supported tidy param values
const (
	TidyTrue  Tidy = "true"
	TidyFalse Tidy = "false"
)

// String returns Tidy as a string
func (t Tidy) String() string {
	return string(t)
}

// Annotations represents a annotations param for osrm5 request
type Annotations string

// Supported annotations param values
const (
	AnnotationsTrue        Annotations = "true"
	AnnotationsFalse       Annotations = "false"
	AnnotationsNodes       Annotations = "nodes"
	AnnotationsDistance    Annotations = "distance"
	AnnotationsDuration    Annotations = "duration"
	AnnotationsDatasources Annotations = "datasources"
	AnnotationsWeight      Annotations = "weight"
	AnnotationsSpeed       Annotations = "speed"
)

// String returns Annotations as a string
func (a Annotations) String() string {
	return string(a)
}

// Steps represents a steps param for osrm5 request
type Steps string

// Supported steps param values
const (
	StepsTrue  Steps = "true"
	StepsFalse Steps = "false"
)

// String returns Steps as a string
func (s Steps) String() string {
	return string(s)
}

// Gaps represents a gaps param for osrm5 match request.
// Allows the input track splitting based on huge timestamp gaps between points.
type Gaps string

// Supported gaps param values
const (
	GapsSplit  Gaps = "split"
	GapsIgnore Gaps = "ignore"
)

// String returns Gaps as a string
func (g Gaps) String() string {
	return string(g)
}

// Geometries represents a geometries param for osrm5
type Geometries string

// Supported geometries param values
const (
	GeometriesPolyline6 Geometries = "polyline6"
	GeometriesGeojson   Geometries = "geojson"
)

// String returns Geometries as a string
func (g Geometries) String() string {
	return string(g)
}

// Overview represents level of overview of geometry in a response
type Overview string

// Available overview values
const (
	OverviewSimplified Overview = "simplified"
	OverviewFull       Overview = "full"
	OverviewFalse      Overview = "false"
)

// String returns Overview as a string
func (o Overview) String() string {
	return string(o)
}

// ContinueStraight represents continue_straight OSRM routing parameter
type ContinueStraight string

// ContinueStraight values
const (
	ContinueStraightDefault ContinueStraight = "default"
	ContinueStraightTrue    ContinueStraight = "true"
	ContinueStraightFalse   ContinueStraight = "false"
)

// String returns ContinueStraight as string
func (c ContinueStraight) String() string {
	return string(c)
}

// Request contains parameters for OSRM query
type Request struct {
	Profile string
	Coords  Coordinates
	Service string
	Options Options
}

// URL generates a url for OSRM request
func (r *Request) URL(serverURL string) (string, error) {
	if r.Service == "" {
		return "", ErrEmptyServiceName
	}
	if r.Profile == "" {
		return "", ErrEmptyProfileName
	}
	if len(r.Coords.MultiPoint) == 0 {
		return "", ErrNoCoordinates
	}
	// http://{server}/{service}/{version}/{profile}/{coordinates}[.{format}]?option=value&option=value
	url := strings.Join([]string{
		serverURL, // server
		r.Service, // service
		"v1",      // version
		r.Profile, // profile
		"polyline(" + url.PathEscape(string(r.Coords.toPolyline())) + ")", // coordinates
	}, "/")
	if len(r.Options) > 0 {
		url += "?" + r.Options.Encode() // options
	}
	return url, nil
}

// Bearing limits the search to segments with given bearing in degrees towards true north in clockwise direction.
type Bearing struct {
	Value, Range uint16
}

func (b Bearing) String() string {
	return fmt.Sprintf("%d,%d", b.Value, b.Range)
}

func Bearings(br []Bearing) string {
	s := make([]string, len(br))
	for i, b := range br {
		s[i] = b.String()
	}
	return strings.Join(s, ";")
}
