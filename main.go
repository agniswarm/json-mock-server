package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/agniswarm/json-mock-server/handlers"
	"github.com/agniswarm/json-mock-server/types"
	"github.com/gin-gonic/gin"
)

// Function to load JSON fixture data with validation
func loadFixture(jsonPath string) ([]types.Route, error) {
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

func main() {
	// Define command-line flags
	filePath := flag.String("file", "", "Path to the JSON file containing routes")
	port := flag.String("port", "3000", "Port on which to run the server (default: 3000)")
	flag.Parse()

	// Check if the file argument is provided
	if *filePath == "" {
		fmt.Println("Error: --file argument is required")
		os.Exit(1)
	}

	// Load the fixture from the specified file
	routes, err := loadFixture(*filePath)
	if err != nil {
		fmt.Printf("Error loading fixture: %v\n", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	// Create a new Gin router
	server := gin.Default()
	// Start the server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: server,
	}

	// Register routes
	handlers.RegisterRoutes(server, routes)

	stop := make(chan os.Signal, 1)

	// Define the exit server route
	server.GET("/exit-server", handlers.ExitServerHandler(httpServer, stop))

	// Handle system signals to shut down gracefully
	signal.Notify(stop, os.Interrupt)

	go func() {
		// Start the server
		log.Printf("Starting server on port %s...\n", *port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %s: %v", *port, err)
		}
	}()

	// Wait for a signal to stop the server
	<-stop
	log.Println("Shutting down server gracefully...")
	if err := httpServer.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
	os.Exit(0)
}
