package types

import (
	"encoding/json"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var geometry = NewGeometryFromMultiPoint(
	orb.MultiPoint{
		{-73.990185, 40.714701},
		{-73.991801, 40.717571},
		{-73.985751, 40.715651},
	},
)

func TestUnmarshalGeometryFromGeojson(t *testing.T) {
	var g Geometry
	in := []byte(`{"type": "LineString", "coordinates": [[-73.982253,40.742926],[-73.985253,40.742926]]}`)

	err := json.Unmarshal(in, &g)

	require.Nil(t, err)
	require.Len(t, g.LineString, 2)
	require.Equal(t, orb.Point{-73.982253, 40.742926}, g.LineString[0])
	require.Equal(t, orb.Point{-73.985253, 40.742926}, g.LineString[1])
}

func TestUnmarshalGeometryFromPolyline(t *testing.T) {
	var g Geometry
	in := []byte(`"nvnalCui}okAkgpk@u}hQf}_l@mbpL"`)

	err := json.Unmarshal(in, &g)

	require.NoError(t, err)
	require.Len(t, g.LineString, 3)
	require.Equal(t, orb.Point{40.123563, -73.965432}, g.LineString[0])
	require.Equal(t, orb.Point{40.423574, -73.235698}, g.LineString[1])
	require.Equal(t, orb.Point{40.645325, -73.973462}, g.LineString[2])
}

func TestUnmarshalGeometryFromNull(t *testing.T) {
	var g Geometry
	in := []byte(`null`)
	err := json.Unmarshal(in, &g)

	require.NoError(t, err)
	require.Equal(t, 0, len(g.LineString))
}

func TestUnmarshalGeometryFromEmptyJSON(t *testing.T) {
	var g Geometry
	in := []byte(`{}`)
	err := json.Unmarshal(in, &g)

	require.Error(t, err)
}

func TestPolylineGeometry(t *testing.T) {
	g := Geometry{
		LineString: geojson.LineString{
			{40.123563, -73.965432},
			{40.423574, -73.235698},
			{40.645325, -73.973462},
		},
	}

	bytes, err := json.Marshal(g)
	require.NoError(t, err)

	assert.Equal(t, `"nvnalCui}okAkgpk@u}hQf}_l@mbpL"`, string(bytes))
}

func TestRequestURLWithEmptyOptions(t *testing.T) {
	req := Request{
		Profile: "something",
		Coords:  geometry,
		Service: "foobar",
	}
	url, err := req.URL("localhost")
	require.Nil(t, err)
	assert.Equal(t, "localhost/foobar/v1/something/polyline(%7BaowFrerbM%7DPbI~Jyd@)", url)
}

func TestRequestURLWithOptions(t *testing.T) {
	opts := Options{}
	opts.Set("baz", "quux")
	req := Request{
		Profile: "something",
		Coords:  geometry,
		Service: "foobar",
		Options: opts,
	}
	url, err := req.URL("localhost")
	require.Nil(t, err)
	assert.Equal(t, "localhost/foobar/v1/something/polyline(%7BaowFrerbM%7DPbI~Jyd@)?baz=quux", url)
}

func TestRequestURLWithEmptyService(t *testing.T) {
	req := Request{}
	url, err := req.URL("localhost")
	require.NotNil(t, err)
	assert.Equal(t, ErrEmptyServiceName, err)
	assert.Empty(t, url)
}

func TestRequestURLWithEmptyProfile(t *testing.T) {
	req := Request{
		Service: "foobar",
	}
	url, err := req.URL("localhost")
	require.NotNil(t, err)
	assert.Equal(t, ErrEmptyProfileName, err)
	assert.Empty(t, url)
}

func TestRequestURLWithoutCoords(t *testing.T) {
	req := Request{
		Profile: "something",
		Service: "foobar",
	}
	url, err := req.URL("localhost")
	require.NotNil(t, err)
	assert.Equal(t, ErrNoCoordinates, err)
	assert.Empty(t, url)
}
