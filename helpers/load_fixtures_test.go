package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFixture(t *testing.T) {
	t.Run("Valid fixture file", func(t *testing.T) {
		filePath := "testdata/valid_fixture.json"
		absPath, _ := filepath.Abs(filePath)

		routes, err := LoadFixture(absPath)
		assert.NoError(t, err)
		assert.Len(t, routes, 2)
		assert.Equal(t, "/test-get", routes[0].Path)
		assert.Equal(t, "/test-post", routes[1].Path)
	})

	t.Run("Invalid JSON content", func(t *testing.T) {
		filePath := "testdata/invalid_fixture.json"
		absPath, _ := filepath.Abs(filePath)

		_, err := LoadFixture(absPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing json")
	})

	t.Run("Non-existent file", func(t *testing.T) {
		_, err := LoadFixture("invalid/path.json")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")
	})
}

func TestMain(m *testing.M) {
	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")

	validFixture := `{
		"routes": [
			{
				"method": "GET",
				"path": "/test-get",
				"status_code": 200,
				"data": "{\"message\": \"hello world\"}"
			},
			{
				"method": "POST",
				"path": "/test-post",
				"status_code": 201,
				"data": "{\"status\": \"created\"}"
			}
		]
	}`
	invalidFixture := `invalid json content`
	invalidRouteFixture := `{
		"routes": [
			{
				"method": "INVALID",
				"path": "/test-invalid",
				"status_code": 200,
				"data": "{\"message\": \"invalid\"}"
			}
		]
	}`

	os.WriteFile("testdata/valid_fixture.json", []byte(validFixture), 0644)
	os.WriteFile("testdata/invalid_fixture.json", []byte(invalidFixture), 0644)
	os.WriteFile("testdata/invalid_route_fixture.json", []byte(invalidRouteFixture), 0644)

	m.Run()
}
