package match

import (
	"github.com/openmarketplaceengine/go-osrm/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyMatchRequestOptions(t *testing.T) {
	cases := []struct {
		name        string
		request     Request
		expectedURI string
	}{
		{
			name:        "empty",
			expectedURI: "geometries=polyline6",
		},
		{
			name: "with timestamps and radiuses",
			request: Request{
				Timestamps: []int64{0, 1, 2},
				Radii:      []float64{0.123123, 0.12312},
			},
			expectedURI: "geometries=polyline6&radiuses=0.123123;0.12312&timestamps=0;1;2",
		},
		{
			name: "with gaps and tidy",
			request: Request{
				Timestamps: []int64{0, 1, 2},
				Radii:      []float64{0.123123, 0.12312},
				Gaps:       types.GapsSplit,
				Tidy:       types.TidyTrue,
			},
			expectedURI: "gaps=split&geometries=polyline6&radiuses=0.123123;0.12312&tidy=true&timestamps=0;1;2",
		},
		{
			name: "with hints",
			request: Request{
				Hints: []string{"a", "b", "c", "d"},
			},
			expectedURI: "geometries=polyline6&hints=a;b;c;d",
		},
		{
			name: "with bearings",
			request: Request{
				Bearings: []types.Bearing{
					{0, 20}, {10, 20},
				},
			},
			expectedURI: "bearings=0%2C20%3B10%2C20&geometries=polyline6",
		},
		{
			name: "custom overview option",
			request: Request{
				Overview:    types.OverviewSimplified,
				Geometries:  types.GeometriesGeojson,
				Annotations: types.AnnotationsFalse,
				Tidy:        types.TidyFalse,
				Steps:       types.StepsFalse,
			},
			expectedURI: "annotations=false&geometries=geojson&overview=simplified&steps=false&tidy=false",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expectedURI, c.request.Request().Options.Encode())
		})
	}
}
