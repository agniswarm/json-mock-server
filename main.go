package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/agniswarm/json-mock-server/handlers"
	"github.com/agniswarm/json-mock-server/helpers"
	"github.com/agniswarm/json-mock-server/notifier"
	"github.com/gin-gonic/gin"
)

func main() {
	// Define command-line flags
	filePath := flag.String("file", "", "Path to the JSON file containing routes")
	port := flag.String("port", "3000", "Port on which to run the server (default: 3000)")
	devMode := flag.Bool("devmode", false, "Enable dev mode to watch for file changes")

	flag.Parse()

	// Check if the file argument is provided
	if *filePath == "" {
		fmt.Println("Error: --file argument is required")
		os.Exit(1)
	}

	// Load the fixture from the specified file
	routes, err := helpers.LoadFixture(*filePath)

	if err != nil {
		fmt.Printf("Error loading routes: %v\n", err)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%s", *port),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	if !*devMode {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	reload := make(chan bool)

	// Activate watch mode only if --devmode flag is passed
	if *devMode {
		go notifier.WatchFileChanges(*filePath, reload)

		go func() {
			for {
				<-reload
				log.Println("Reloading routes...")
				routes, err := helpers.LoadFixture(*filePath)
				if err != nil {
					log.Printf("Error reloading routes: %v", err)
					continue
				}

				server = gin.New() // Clear existing routes

				httpServer.Shutdown(context.TODO())

				httpServer = &http.Server{
					Addr: fmt.Sprintf(":%s", *port),
				}

				if err := handlers.CheckDuplicateRoutes(routes); err != nil {
					log.Println(err.Error())
					log.Println("Server will not start, please duplicate route error")
					continue
				}

				httpServer.Handler = server
				go helpers.StartServer(httpServer, routes, server, stop)
			}
		}()
	}

	go helpers.StartServer(httpServer, routes, server, stop)

	// Wait for a signal to stop the server
	<-stop
	log.Println("Shutting down server gracefully...")
	if err := httpServer.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
	os.Exit(0)
}
