package route

import (
	"github.com/openmarketplaceengine/go-osrm/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyRouteRequestOptions(t *testing.T) {
	req := Request{}
	assert.Equal(
		t,
		"geometries=polyline6",
		req.Request().Options.Encode())
}

func TestRouteRequestOptionsWithBearings(t *testing.T) {
	req := Request{
		Bearings: []types.Bearing{
			{60, 380},
			{45, 180},
		},
		ContinueStraight: types.ContinueStraightTrue,
	}
	assert.Equal(
		t,
		"bearings=60%2C380%3B45%2C180&continue_straight=true&geometries=polyline6",
		req.Request().Options.Encode())
}

func TestRouteRequestOverviewOption(t *testing.T) {
	req := Request{
		Overview:         types.OverviewFull,
		ContinueStraight: types.ContinueStraightTrue,
	}
	assert.Equal(
		t,
		"continue_straight=true&geometries=polyline6&overview=full",
		req.Request().Options.Encode())
}

func TestRouteRequestGeometryOption(t *testing.T) {
	req := Request{
		Geometries:       types.GeometriesPolyline6,
		Annotations:      types.AnnotationsFalse,
		Steps:            types.StepsFalse,
		ContinueStraight: types.ContinueStraightTrue,
	}
	assert.Equal(
		t,
		"annotations=false&continue_straight=true&geometries=polyline6&steps=false",
		req.Request().Options.Encode())
}
