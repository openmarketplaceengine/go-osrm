package osrm_test

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/go-osrm/route"
	"github.com/openmarketplaceengine/go-osrm/types"
	"log"

	osrm "github.com/openmarketplaceengine/go-osrm"
	geo "github.com/paulmach/orb"
)

func ExampleOSRM_Route() {
	client := osrm.NewFromURL("https://router.project-osrm.org")

	resp, err := client.Route(context.Background(), route.Request{
		Profile: "car",
		Coordinates: types.NewGeometryFromMultiPoint(geo.MultiPoint{
			{-73.87946, 40.75833},
			{-73.87925, 40.75837},
			{-73.87918, 40.75837},
			{-73.87911, 40.75838},
		}),
		Steps:       types.StepsTrue,
		Annotations: types.AnnotationsTrue,
		Overview:    types.OverviewFalse,
		Geometries:  types.GeometriesPolyline6,
	})
	if err != nil {
		log.Fatalf("route failed: %v", err)
	}

	fmt.Println(len(resp.Routes))

	// Output:
	// 1
}
