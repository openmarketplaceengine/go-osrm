package types

import (
	"encoding/json"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var geometry = geojson.LineString{
	{-73.990185, 40.714701},
	{-73.991801, 40.717571},
	{-73.985751, 40.715651},
}

func Test_Coordinates_toPolyline(t *testing.T) {
	c := coordinatesFromLineString(geometry)
	require.Equal(t, ``, string(c.toPolyline()))
}

func TestUnmarshal(t *testing.T) {
	tests := map[string]struct {
		b      []byte
		expect func(t *testing.T, c *Coordinates, err error)
	}{
		"From GeoJSON": {
			b: []byte(`{"type": "LineString", "coordinates": [[-73.982253,40.742926],[-73.985253,40.742926]]}`),
			expect: func(t *testing.T, c *Coordinates, err error) {
				require.NoError(t, err)
				p := c.MultiPoint
				require.Len(t, p, 2)
				require.Equal(t, orb.Point{-73.982253, 40.742926}, p[0])
				require.Equal(t, orb.Point{-73.985253, 40.742926}, p[1])
			},
		},
		"From Polyline": {
			b: []byte(`"nvnalCui}okAkgpk@u}hQf}_l@mbpL"`),
			expect: func(t *testing.T, c *Coordinates, err error) {
				require.NoError(t, err)
				p := c.MultiPoint
				require.Len(t, p, 3)
				require.Equal(t, orb.Point{40.123563, -73.965432}, p[0])
				require.Equal(t, orb.Point{40.423574, -73.235698}, p[1])
				require.Equal(t, orb.Point{40.645325, -73.973462}, p[2])
			},
		},
		"Null": {
			b: []byte(`null`),
			expect: func(t *testing.T, c *Coordinates, err error) {
				require.NoError(t, err)
				require.Equal(t, 0, len(c.MultiPoint))
			},
		},
		"Empty": {
			b: []byte(`{}`),
			expect: func(t *testing.T, c *Coordinates, err error) {
				require.Error(t, err)
				require.EqualError(t, err, `foobar`)
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			var c Coordinates
			err := json.Unmarshal(tc.b, &c)
			tc.expect(t, &c, err)
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := map[string]struct {
		coordinates Coordinates
		expect      func(t *testing.T, b []byte, err error)
	}{
		"OK": {
			coordinates: coordinatesFromLineString(
				geojson.LineString{
					{40.123563, -73.965432},
					{40.423574, -73.235698},
					{40.645325, -73.973462},
				},
			),
			expect: func(t *testing.T, b []byte, err error) {
				require.NoError(t, err)
				assert.Equal(t, `"nvnalCui}okAkgpk@u}hQf}_l@mbpL"`, string(b))
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			b, err := json.Marshal(tc.coordinates)
			tc.expect(t, b, err)
		})
	}
}

func TestBuildURL(t *testing.T) {
	tests := map[string]struct {
		buildRequest func() Request
		expect       func(t *testing.T, actual string, err error)
	}{
		"Empty Options": {
			buildRequest: func() Request {
				return Request{
					Profile: "something",
					Coords:  coordinatesFromLineString(geometry),
					Service: "foobar",
				}
			},
			expect: func(t *testing.T, actual string, err error) {
				require.NoError(t, err)
				assert.Equal(t, "localhost/foobar/v1/something/polyline(%7BaowFrerbM%7DPbI~Jyd@)", actual)
			},
		},
		"With Options": {
			buildRequest: func() Request {
				opts := Options{}
				opts.Set("baz", "quux")
				return Request{
					Profile: "something",
					Coords:  coordinatesFromLineString(geometry),
					Service: "foobar",
					Options: opts,
				}
			},
			expect: func(t *testing.T, actual string, err error) {
				require.NoError(t, err)
				assert.Equal(t, "localhost/foobar/v1/something/polyline(%7BaowFrerbM%7DPbI~Jyd@)?baz=quux", actual)
			},
		},
		"With Empty Service": {
			buildRequest: func() Request {
				return Request{}
			},
			expect: func(t *testing.T, actual string, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrEmptyServiceName, err)
				assert.Empty(t, actual)
			},
		},
		"With Empty Profile": {
			buildRequest: func() Request {
				return Request{
					Service: "foobar",
				}
			},
			expect: func(t *testing.T, actual string, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrEmptyProfileName, err)
				assert.Empty(t, actual)
			},
		},
		"Without Coords": {
			buildRequest: func() Request {
				return Request{
					Profile: "something",
					Service: "foobar",
				}
			},
			expect: func(t *testing.T, actual string, err error) {
				require.Error(t, err)
				assert.Equal(t, ErrNoCoordinates, err)
				assert.Empty(t, actual)
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			req := tc.buildRequest()
			url, err := req.URL("localhost")
			tc.expect(t, url, err)
		})
	}
}
