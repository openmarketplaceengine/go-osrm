package types

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/twpayne/go-polyline"
)

type Coordinates struct {
	geojson.MultiPoint
}

func (c Coordinates) MarshalJSON() ([]byte, error) {
	return c.toPolyline(), nil
}

func (c *Coordinates) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Is it a string?
	var encoded string
	if err := json.Unmarshal(data, &encoded); err == nil {
		codec := polyline.Codec{Dim: 2, Scale: 1e6}
		coords, _, err := codec.DecodeCoords([]byte(encoded))
		if err != nil {
			return err
		}
		c.MultiPoint = coordsToMultiPoint(coords)
		return nil
	}

	geom, err := geojson.UnmarshalGeometry(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal geojson geometry, err: %v", err)
	}
	if geom.Type != "LineString" {
		return fmt.Errorf("unexpected geometry type: %v", geom.Type)
	}

	var mp geojson.MultiPoint
	for _, p := range geom.Coordinates.(orb.LineString) {
		mp = append(mp, orb.Point([2]float64{p.X(), p.Y()}))
	}
	c.MultiPoint = mp

	return nil
}

func (c Coordinates) toPolyline(in ...float64) []byte {
	scale := 1e6
	if len(in) > 0 {
		scale = in[0]
	}
	var coords [][]float64
	for _, p := range c.MultiPoint {
		coords = append(coords, []float64{p.X(), p.Y()})
	}
	codec := polyline.Codec{Dim: 2, Scale: scale}
	return codec.EncodeCoords(nil, coords)
}

func coordinatesFromLineString(ls geojson.LineString) Coordinates {
	return Coordinates{geojson.MultiPoint(ls)}
}

func coordsToMultiPoint(coords [][]float64) geojson.MultiPoint {
	var mp geojson.MultiPoint
	for _, p := range coords {
		mp = append(mp, orb.Point([2]float64{p[1], p[0]}))
	}
	return mp
}
