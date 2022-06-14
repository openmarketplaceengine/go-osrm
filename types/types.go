package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"math"
	"net/url"
	"strings"

	"github.com/paulmach/orb"
)

const (
	polyline5Factor = 1.0e5
	polyline6Factor = 1.0e6
)

// Geometry represents a points set
type Geometry struct {
	geojson.LineString
}

// NewGeometryFromPath creates a geometry from a path.
func NewGeometryFromPath(path geojson.LineString) Geometry {
	return Geometry{LineString: path}
}

// NewGeometryFromMultiPoint creates a geometry from points set.
func NewGeometryFromMultiPoint(ps orb.MultiPoint) Geometry {
	return NewGeometryFromPath(geojson.LineString(ps))
}

// Polyline generates a polyline in Google format
func (g *Geometry) Polyline(factor ...int) string {
	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	var pLat int
	var pLng int

	var result bytes.Buffer
	scratch1 := make([]byte, 0, 50)
	scratch2 := make([]byte, 0, 50)

	for _, p := range g.LineString {
		lat5 := int(math.Floor(p.Lat()*f + 0.5))
		lng5 := int(math.Floor(p.Lon()*f + 0.5))

		deltaLat := lat5 - pLat
		deltaLng := lng5 - pLng

		pLat = lat5
		pLng = lng5

		result.Write(append(encodeSignedNumber(deltaLat, scratch1), encodeSignedNumber(deltaLng, scratch2)...))

		scratch1 = scratch1[:0]
		scratch2 = scratch2[:0]
	}

	return result.String()
}

func encodeSignedNumber(num int, result []byte) []byte {
	shiftedNum := num << 1

	if num < 0 {
		shiftedNum = ^shiftedNum
	}

	for shiftedNum >= 0x20 {
		result = append(result, byte(0x20|(shiftedNum&0x1f)+63))
		shiftedNum >>= 5
	}

	return append(result, byte(shiftedNum+63))
}

// UnmarshalJSON parses a geo path from points set or a polyline
func (g *Geometry) UnmarshalJSON(b []byte) error {
	return g.LineString.UnmarshalJSON(b)
}

// MarshalJSON generates a polyline in Google polyline6 format
func (g Geometry) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.Polyline(polyline6Factor))
}

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
	Coords  Geometry
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
	if len(r.Coords.LineString) == 0 {
		return "", ErrNoCoordinates
	}
	// http://{server}/{service}/{version}/{profile}/{coordinates}[.{format}]?option=value&option=value
	url := strings.Join([]string{
		serverURL, // server
		r.Service, // service
		"v1",      // version
		r.Profile, // profile
		"polyline(" + url.PathEscape(r.Coords.Polyline(polyline5Factor)) + ")", // coordinates
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
