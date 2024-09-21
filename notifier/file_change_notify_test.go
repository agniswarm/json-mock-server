package notifier

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatchFileChanges(t *testing.T) {
	testFilePath := "testdata/test_file.json"
	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")

	reload := make(chan bool)

	err := os.WriteFile(testFilePath, []byte(`{"key": "value"}`), 0644)
	assert.NoError(t, err)

	go WatchFileChanges(testFilePath, reload)

	// Give some time for the watcher to start
	time.Sleep(1 * time.Second)

	// Modify the file
	err = os.WriteFile(testFilePath, []byte(`{"key1": "new value1"}`), 0644)
	assert.NoError(t, err)

	select {
	case <-reload:
		log.Println("File modification detected")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for file modification event")
	}
}
