package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/agniswarm/json-mock-server/types"
)

// Function to load JSON fixture data with validation
func LoadFixture(jsonPath string) ([]types.Route, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return []types.Route{}, fmt.Errorf("failed to read file: %v", err)
	}

	var fixture types.Fixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return []types.Route{}, fmt.Errorf("error parsing json: %v", err)
	}

	// Validate routes and their data
	for _, route := range fixture.Routes {
		if err := route.ValidateRoute(); err != nil {
			return []types.Route{}, err
		}
	}

	return fixture.Routes, nil
}
