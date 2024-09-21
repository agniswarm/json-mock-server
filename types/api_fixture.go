package types

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Fixture struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Path       string      `json:"path"`
	Method     string      `json:"method"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"status_code"`
}

// Function to validate a route's data
func (route Route) ValidateRoute() error {
	switch data := route.Data.(type) {
	case string:
		if strings.HasPrefix(data, "json://") {
			filePath := strings.TrimPrefix(data, "json://")
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return fmt.Errorf("invalid file path for route %s: %v", route.Path, err)
			}
			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				return fmt.Errorf("data file does not exist for route %s: %s", route.Path, absPath)
			}

			// Check if the file is valid JSON
			if err := validateJSONFile(absPath); err != nil {
				return fmt.Errorf("invalid JSON in data file for route %s: %v", route.Path, err)
			}
		}
	case map[string]interface{}, []interface{}:
		// Valid JSON structure
		return nil
	default:
		return fmt.Errorf("invalid data for route %s: must be a string or valid JSON structure", route.Path)
	}
	return nil
}

// Function to validate the JSON structure of a file
func validateJSONFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read data file: %v", err)
	}

	var content interface{}
	return json.Unmarshal(data, &content)
}
