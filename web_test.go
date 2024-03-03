package gitearelease

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchData(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the path part of the URL
		if r.URL.Path != "/data" {
			t.Errorf("Expected URL path: /data, got: %s", r.URL.Path)
		}

		switch r.URL.Path {
		case "/data":
			// Set the response status code
			w.WriteHeader(http.StatusOK)

			// Set the response body
			_, err := w.Write([]byte("Test data"))
			if err != nil {
				t.Errorf("Failed to write the response: %v", err)
			}
		case "/error":
			// Simulate an error response
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	data, err := fetchData(server.URL + "/data")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the returned data
	expectedData := "Test data"
	if string(data) != expectedData {
		t.Errorf("Expected data: %s, got: %s", expectedData, data)
	}

	// Test the error case by requesting the error path
	_, err = fetchData(server.URL + "/error")
	if err == nil {
		t.Errorf("Expected an error but got nil")
	}
}
