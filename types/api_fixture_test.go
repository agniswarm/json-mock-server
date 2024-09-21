package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute_ValidateRoute(t *testing.T) {
	// Setup test data directory
	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")

	validJSONPath := "testdata/valid.json"
	invalidJSONPath := "testdata/invalid.json"

	os.WriteFile(validJSONPath, []byte(`{"key": "value"}`), 0644)
	os.WriteFile(invalidJSONPath, []byte(`invalid json content`), 0644)

	tests := []struct {
		name    string
		route   Route
		wantErr bool
	}{
		{
			name: "Valid JSON string data",
			route: Route{
				Path:       "/test",
				Method:     "GET",
				Data:       `{"key": "value"}`,
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "Valid JSON file data",
			route: Route{
				Path:       "/test",
				Method:     "GET",
				Data:       "json://" + validJSONPath,
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "Invalid JSON file data",
			route: Route{
				Path:       "/test",
				Method:     "GET",
				Data:       "json://" + invalidJSONPath,
				StatusCode: 200,
			},
			wantErr: true,
		},
		{
			name: "Non-existent JSON file data",
			route: Route{
				Path:       "/test",
				Method:     "GET",
				Data:       "json://nonexistent.json",
				StatusCode: 200,
			},
			wantErr: true,
		},
		{
			name: "Invalid data type",
			route: Route{
				Path:       "/test",
				Method:     "GET",
				Data:       12345,
				StatusCode: 200,
			},
			wantErr: true,
		},
		{
			name: "Valid JSON structure data",
			route: Route{
				Path:       "/test",
				Method:     "GET",
				Data:       map[string]interface{}{"key": "value"},
				StatusCode: 200,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.route.ValidateRoute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateJSONFile(t *testing.T) {
	// Setup test data directory
	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")

	validJSONPath := "testdata/valid.json"
	invalidJSONPath := "testdata/invalid.json"

	os.WriteFile(validJSONPath, []byte(`{"key": "value"}`), 0644)
	os.WriteFile(invalidJSONPath, []byte(`invalid json content`), 0644)

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Valid JSON file",
			filePath: validJSONPath,
			wantErr:  false,
		},
		{
			name:     "Invalid JSON file",
			filePath: invalidJSONPath,
			wantErr:  true,
		},
		{
			name:     "Non-existent file",
			filePath: "testdata/nonexistent.json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONFile(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
