package types

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Coordinates_toPolyline(t *testing.T) {
	c := coordinatesFromLineString(geojson.LineString{
		{-73.990185, 40.714701},
	})
	b := c.toPolyline()
	s := "\"" + string(b) + "\""
	assert.Equal(t, `"pa_clCy{_tlA"`, s)
	var c2 Coordinates
	err := json.Unmarshal([]byte(s), &c2)
	require.NoError(t, err)
	b, err = json.Marshal(c2.MultiPoint)
	require.NoError(t, err)
	fmt.Println(string(b))
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
				require.EqualError(t, err, `failed to unmarshal geojson geometry, err: geojson: invalid geometry`)
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