package helpers

import (
	"log"
	"net/http"
	"os"

	"github.com/agniswarm/json-mock-server/handlers"
	"github.com/agniswarm/json-mock-server/types"
	"github.com/gin-gonic/gin"
)

func StartServer(httpServer *http.Server, routes []types.Route, server *gin.Engine, stop chan os.Signal) {

	httpServer.Handler = server

	// Register routes
	if err := handlers.RegisterRoutes(server, routes); err != nil {
		log.Fatalf("%v", err)
	}

	// Define the exit server route
	server.GET("/exit-server", handlers.ExitServerHandler(httpServer, stop))

	// Handle system signals to shut down gracefully

	// Start the server
	log.Printf("Starting server on port %s...\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on port %s: %v", httpServer.Addr, err)
	}
}
