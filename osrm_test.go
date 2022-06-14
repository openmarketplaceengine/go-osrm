package osrm

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/go-osrm/match"
	"github.com/openmarketplaceengine/go-osrm/nearest"
	"github.com/openmarketplaceengine/go-osrm/route"
	"github.com/openmarketplaceengine/go-osrm/table"
	"github.com/openmarketplaceengine/go-osrm/types"
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var geometry = types.NewGeometryFromMultiPoint(
	orb.MultiPoint{
		{-73.990185, 40.714701},
		{-73.991801, 40.717571},
		{-73.985751, 40.715651},
	},
)

func fixturedJSON(name string) []byte {
	data, err := ioutil.ReadFile("testdata/" + name + ".json")
	if err != nil {
		log.Fatalf("osrm5: failed to load a fixture %s, err: %s", name, err)
	}
	return data
}

func fixturedHTTPHandler(name string, assertURL func(path, query string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assertURL(r.URL.Path, r.URL.RawQuery)
		_, _ = fmt.Fprintln(w, string(fixturedJSON(name)))
	}
}

func TestDefaultOSRMConfig(t *testing.T) {
	osrm := New()

	assert.Equal(t, defaultServerURL, osrm.client.serverURL)
}

func TestErrorWithTimeout(t *testing.T) {
	osrm := NewFromURLWithTimeout("http://25.0.0.1", 500*time.Microsecond)

	var nothing response

	req := types.Request{
		Service: "nothing",
		Profile: "nothing",
		Coords:  geometry,
	}

	err := osrm.query(context.Background(), &req, nothing)
	require.Error(t, err)
}

func TestErrorOnRequest(t *testing.T) {
	ts := httptest.NewServer(fixturedHTTPHandler("invalid_query_response", func(path, query string) {}))
	defer ts.Close()

	osrm := NewFromURL(ts.URL)

	geom := types.NewGeometryFromMultiPoint(orb.MultiPoint{{0.1, 0.1}})

	assert := func(t *testing.T, err error) {
		t.Helper()
		require.EqualError(t, err, "InvalidQuery - Query string malformed close to position 28")
		assert.Equal(t, types.ErrorCodeInvalidQuery, err.(types.ResponseStatus).ErrCode())
	}

	t.Run("route", func(t *testing.T) {
		_, err := osrm.Route(context.Background(), route.Request{
			Profile:     "car",
			Coordinates: geom,
		})

		assert(t, err)
	})

	t.Run("match", func(t *testing.T) {
		_, err := osrm.Match(context.Background(), match.Request{
			Profile:     "car",
			Coordinates: geom,
		})

		assert(t, err)
	})

	t.Run("table", func(t *testing.T) {
		_, err := osrm.Table(context.Background(), table.Request{
			Profile:     "car",
			Coordinates: geom,
		})

		assert(t, err)
	})

	t.Run("nearest", func(t *testing.T) {
		_, err := osrm.Nearest(context.Background(), nearest.Request{
			Profile:     "car",
			Coordinates: geom,
		})

		assert(t, err)
	})
}

func TestRouteRequest(t *testing.T) {
	ts := httptest.NewServer(fixturedHTTPHandler("route_response_full", func(path, query string) {
		assert.Equal(t, "/route/v1/car/polyline({aowFrerbM}PbI~Jyd@)", path)
		assert.Equal(t, "annotations=true&continue_straight=true&geometries=polyline6&overview=full", query)
	}))
	defer ts.Close()

	osrm := NewFromURL(ts.URL)

	r, err := osrm.Route(context.Background(), route.Request{
		Profile:          "car",
		Coordinates:      geometry,
		Annotations:      types.AnnotationsTrue,
		Geometries:       types.GeometriesPolyline6,
		Overview:         types.OverviewFull,
		ContinueStraight: types.ContinueStraightTrue,
	})

	require := require.New(t)

	require.NoError(err)
	require.NotNil(r)

	// response
	require.Equal("2017-11-17T21:43:02Z", r.DataVersion)
	// routes
	require.Len(r.Routes, 1)
	route := r.Routes[0]
	require.Equal(float32(1190.5), route.Distance)
	require.Equal(float32(92.2), route.Duration)
	// routes/legs
	require.Len(route.Legs, 2)
	leg0 := route.Legs[0]
	require.Equal(float32(637.5), leg0.Distance)
	require.Equal(float32(58.0), leg0.Duration)
	// routes/annotations
	annotation := leg0.Annotation
	require.Len(annotation.Duration, 14)
	require.Len(annotation.Distance, 14)
	// routes/legs/steps
	require.Len(leg0.Steps, 7)
	// routes/legs/steps[0]
	step0 := leg0.Steps[0]
	require.Equal("driving", step0.Mode)
	require.Equal("", step0.Name)
	require.Equal(float32(5.0), step0.Duration)
	require.Equal(float32(33.1), step0.Distance)
	require.Equal(types.Geometry{
		LineString: geojson.LineString{
			{-73.9902, 40.7147},
			{-73.99023, 40.7146},
			{-73.99025, 40.71441},
		},
	}, step0.Geometry)
}

func TestTableRequest(t *testing.T) {
	ts := httptest.NewServer(fixturedHTTPHandler("table_response_full", func(path, query string) {
		assert.Equal(t, "/table/v1/car/polyline({aowFrerbM}PbI~Jyd@)", path)
		assert.Empty(t, query)
	}))
	defer ts.Close()

	osrm := NewFromURL(ts.URL)

	r, err := osrm.Table(context.Background(), table.Request{Profile: "car", Coordinates: geometry})

	require := require.New(t)

	require.NoError(err)
	require.NotNil(r)

	require.Len(r.Durations, 3)
	require.Equal([]float32{0, 39, 46.8}, r.Durations[0])
	require.Equal([]float32{39.5, 0, 34.2}, r.Durations[1])
	require.Equal([]float32{47.2, 34.2, 0}, r.Durations[2])
}

func TestMatchRequest(t *testing.T) {
	ts := httptest.NewServer(fixturedHTTPHandler("match_response_full", func(path, query string) {
		assert.Equal(t, "/match/v1/car/polyline({aowFrerbM}PbI~Jyd@)", path)
		assert.Equal(t, "geometries=polyline6", query)
	}))
	defer ts.Close()

	osrm := NewFromURL(ts.URL)

	r, err := osrm.Match(context.Background(), match.Request{
		Profile:     "car",
		Coordinates: geometry,
	})

	require := require.New(t)

	require.NoError(err)
	require.NotNil(r)

	// response
	require.Equal("new", r.DataVersion)
	// matchings
	require.Len(r.Matchings, 1)
	matching := r.Matchings[0]
	require.Equal(0.023898, matching.Confidence)
	require.Equal(float32(1035.3), matching.Distance)
	require.Equal(float32(79.0), matching.Duration)
	// matchings/legs
	require.Len(matching.Legs, 2)
	require.Len(matching.Legs[0].Annotation.Nodes, 11)
	require.Len(matching.Legs[1].Annotation.Nodes, 15)
}

func TestNearestRequest(t *testing.T) {
	ts := httptest.NewServer(fixturedHTTPHandler("nearest_response_full", func(path, query string) {
		assert.Equal(t, "/nearest/v1/car/polyline(edswF|`sbM)", path)
		assert.Equal(t, "number=5", query)
	}))
	defer ts.Close()

	osrm := NewFromURL(ts.URL)

	r, err := osrm.Nearest(context.Background(), nearest.Request{
		Profile: "car",
		Coordinates: types.NewGeometryFromMultiPoint(orb.MultiPoint{
			{-73.994550, 40.735551},
		}),
		Number: 5,
	})

	require := require.New(t)

	require.NoError(err)
	require.NotNil(r)

	assert.Len(t, r.Waypoints, 5)
	assert.Equal(t, "XRAFgP___3-SAAAAXAEAAAAAAAAAAAAAfCMzQjwphUMAAAAAAAAAAJIAAABcAQAAAAAAAAAAAACVCQAAqu6W-xSTbQLK7pb7P5NtAgAAvxLg85BF", r.Waypoints[0].Hint)
	assert.Equal(t, "-B4FgP___3_mAAAAZAAAAAAAAAAAAAAAcKqmQmexskAAAAAAAAAAAOYAAABkAAAAAAAAAAAAAACVCQAAmvCW-32SbQLK7pb7P5NtAgAADw3g85BF", r.Waypoints[1].Hint)
	assert.Equal(t, "qRoFgP___38kAwAAyAAAAAAAAAAAAAAAjOUyQwAAAAAAAAAAAAAAACQDAADIAAAAAAAAAAAAAACVCQAAevCW-1GSbQLK7pb7P5NtAgAAvxLg85BF", r.Waypoints[2].Hint)
	assert.Equal(t, "XhAFgP___38AAAAAWwAAAAAAAAAAAAAAAAAAANvEokIAAAAAAAAAAAAAAABbAAAAAAAAAAAAAACVCQAAevCW-1GSbQLK7pb7P5NtAgAADw3g85BF", r.Waypoints[3].Hint)
	assert.Equal(t, "-h4FgJQVyYAyAAAA2AAAAAAAAAAAAAAAU4QzQm0XQUMAAAAAAAAAADIAAADYAAAAAAAAAAAAAACVCQAAp_CW-8-VbQLK7pb7P5NtAgAArxLg85BF", r.Waypoints[4].Hint)
}
