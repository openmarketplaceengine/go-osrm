package nearest

import (
	"github.com/openmarketplaceengine/go-osrm/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNearestRequestOverviewOption(t *testing.T) {
	req := Request{
		Number: 2,
		Bearings: []types.Bearing{
			{60, 380},
		},
	}
	assert.Equal(
		t,
		"bearings=60%2C380&number=2",
		req.Request().Options.Encode())

	req = Request{
		Bearings: []types.Bearing{
			{60, 380},
		},
	}
	assert.Equal(
		t,
		"bearings=60%2C380",
		req.Request().Options.Encode())
}
