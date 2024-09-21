package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/agniswarm/json-mock-server/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	routes := []types.Route{
		{
			Method:     http.MethodGet,
			Path:       "/test-get",
			StatusCode: http.StatusOK,
			Data:       `{"message": "hello world"}`,
		},
		{
			Method:     http.MethodPost,
			Path:       "/test-post",
			StatusCode: http.StatusCreated,
			Data:       `{"status": "created"}`,
		},
	}

	RegisterRoutes(router, routes)

	t.Run("GET /test-get", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test-get", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var responseData map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", responseData["message"])
	})

	t.Run("POST /test-post", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/test-post", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		var responseData map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "created", responseData["status"])
	})
}

func TestPrepareResponseData(t *testing.T) {
	t.Run("String data", func(t *testing.T) {
		data := `{"key": "value"}`
		result, err := prepareResponseData(data)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"key": "value"}, result)
	})

	t.Run("File data", func(t *testing.T) {
		filePath := "testdata/test.json"
		absPath, _ := filepath.Abs(filePath)
		data := "json://" + absPath

		result, err := prepareResponseData(data)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"fileKey": "fileValue"}, result)
	})

	t.Run("Invalid JSON string", func(t *testing.T) {
		data := `invalid json`
		result, err := prepareResponseData(data)
		assert.NoError(t, err)
		assert.Equal(t, "invalid json", result)
	})

	t.Run("Unsupported data type", func(t *testing.T) {
		data := 12345
		_, err := prepareResponseData(data)
		assert.Error(t, err)
	})
}

func TestLoadData(t *testing.T) {
	t.Run("Valid file", func(t *testing.T) {
		filePath := "testdata/test.json"
		absPath, _ := filepath.Abs(filePath)

		result, err := loadData(absPath)
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"fileKey": "fileValue"}, result)
	})

	t.Run("Invalid file path", func(t *testing.T) {
		_, err := loadData("invalid/path.json")
		assert.Error(t, err)
	})

	t.Run("Invalid JSON content", func(t *testing.T) {
		filePath := "testdata/invalid.json"
		absPath, _ := filepath.Abs(filePath)

		_, err := loadData(absPath)
		assert.Error(t, err)
	})
}

// Create a testdata directory with test.json and invalid.json files for testing
func TestMain(m *testing.M) {
	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")

	os.WriteFile("testdata/test.json", []byte(`{"fileKey": "fileValue"}`), 0644)
	os.WriteFile("testdata/invalid.json", []byte(`invalid json content`), 0644)

	m.Run()
}
