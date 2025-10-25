package traefik_conditional_headers

import (
	"net/http"
	"net/http/httptest"
)

// createTestRequest creates an HTTP request for testing.
func createTestRequest(host string) *http.Request {
	req := httptest.NewRequest("GET", "https://"+host+"/test", nil)
	req.Host = host
	return req
}

// createTestResponse creates a ResponseRecorder for testing.
func createTestResponse() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// mockNextHandler creates a simple handler that records it was called.
func mockNextHandler() http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusOK)
		_, err := responseWriter.Write([]byte("OK"))
		if err != nil {
			return
		}
	})
}

// headerCaptureHandler creates a handler that captures the request headers.
func headerCaptureHandler(capturedHeadersRef *map[string][]string) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		*capturedHeadersRef = make(map[string][]string)
		for key, values := range request.Header {
			(*capturedHeadersRef)[key] = values
		}
		responseWriter.WriteHeader(http.StatusOK)
		_, err := responseWriter.Write([]byte("OK"))
		if err != nil {
			return
		}
	})
}
