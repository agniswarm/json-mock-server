package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// /exit-server route handler
func ExitServerHandler(server *http.Server, stop chan os.Signal) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Shutting down the server...\n")
		go func() {
			time.Sleep(500 * time.Millisecond) // Give the client some time to see the response
			stop <- os.Interrupt
		}()
	}
}
