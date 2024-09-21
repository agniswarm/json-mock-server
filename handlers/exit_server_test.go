package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	// spell-checker: disable
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// spell-checker: enable
)

type MockOS struct {
	mock.Mock
}

func (m *MockOS) Exit(code int) {
	m.Called(code)
}

func TestExitServerHandlerWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	mockOS := new(MockOS)
	mockOS.On("Exit", 0).Once()

	signalChan := make(chan os.Signal, 1)
	router.GET("/exit-server", ExitServerHandler(server, signalChan))

	t.Run("GET /exit-server", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/exit-server", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "Shutting down the server...\n", resp.Body.String())

		// Wait a bit to ensure the server would have exited
		time.Sleep(1 * time.Second)

	})
}
