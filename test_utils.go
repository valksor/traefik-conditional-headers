package traefik_conditional_headers

import (
	"net/http"
	"net/http/httptest"
)

// createTestRequest creates an HTTP request for testing
func createTestRequest(host string, path string) *http.Request {
	req := httptest.NewRequest("GET", "https://"+host+path, nil)
	req.Host = host
	return req
}

// createTestResponse creates a ResponseRecorder for testing
func createTestResponse() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// mockNextHandler creates a simple handler that records it was called
func mockNextHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte("OK"))
		if err != nil {
			return
		}
	})
}

// headerCaptureHandler creates a handler that captures the request headers
func headerCaptureHandler(capturedHeaders *map[string][]string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		*capturedHeaders = make(map[string][]string)
		for key, values := range req.Header {
			(*capturedHeaders)[key] = values
		}
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte("OK"))
		if err != nil {
			return
		}
	})
}
