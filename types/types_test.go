package types

import (
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
