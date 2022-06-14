package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTableRequestOptions(t *testing.T) {
	req := Request{}
	assert.Empty(t, req.Request().Options.Encode())
}

func TestNotEmptyTableRequestOptions(t *testing.T) {
	req := Request{
		Sources:      []int{0, 1, 2},
		Destinations: []int{1, 3},
	}
	assert.Equal(t, "destinations=1;3&sources=0;1;2", req.Request().Options.Encode())
}
