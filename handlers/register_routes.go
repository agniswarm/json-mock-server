package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/agniswarm/json-mock-server/types"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, routes []types.Route) error {
	for _, route := range routes {
		if err := fillRouteData(server, route); err != nil {
			return err
		}
	}
	return nil
}

// Function to validate duplicate routes if the method and routes are same
func CheckDuplicateRoutes(routes []types.Route) error {
	uniqueRoutes := make(map[string]bool)
	for _, route := range routes {
		routeKey := fmt.Sprintf("%s:%s", route.Method, route.Path)
		if _, ok := uniqueRoutes[routeKey]; ok {
			return fmt.Errorf("duplicate route found: %v %v", route.Method, route.Path)
		}
		uniqueRoutes[routeKey] = true
	}
	return nil
}

func fillRouteData(server *gin.Engine, route types.Route) error {
	if route.StatusCode == 0 {
		route.StatusCode = http.StatusOK
	}

	responseData, err := prepareResponseData(route.Data)
	if err != nil {
		return fmt.Errorf("error preparing response data for route: %v", route.Path)
	}

	okRoute := false

	for _, m := range []string{"GET", "POST"} {
		if m == route.Method {
			okRoute = true
			break
		}
	}

	if !okRoute {
		return fmt.Errorf("invalid route method: %v", route.Method)
	}

	if route.Method == http.MethodGet {
		server.GET(route.Path, func(c *gin.Context) {
			c.JSON(route.StatusCode, responseData)
		})
	} else if route.Method == http.MethodPost {
		server.POST(route.Path, func(c *gin.Context) {
			c.JSON(route.StatusCode, responseData)
		})
	}
	return nil
}

// Function to prepare response data
func prepareResponseData(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case string:
		// Check if the data is a file path with file:// prefix
		if strings.HasPrefix(v, "json://") {
			filePath := strings.TrimPrefix(v, "json://")
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return nil, fmt.Errorf("invalid file path: %v", err)
			}
			return loadData(absPath)
		}
		// If it's a regular string, try to parse it as JSON
		var jsonData interface{}
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			// If parsing fails, return the string as is
			return strings.TrimSpace(v), nil
		}
		return jsonData, nil
	case map[string]interface{}, []interface{}:
		return v, nil // For maps and slices, return as is
	default:
		return nil, fmt.Errorf("unsupported data type: %T", v)
	}
}

// Function to load data from a file
func loadData(filePath string) (interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %v", err)
	}

	var content interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("error parsing data json: %v", err)
	}
	return content, nil
}
