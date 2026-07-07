package traefik_conditional_headers

import (
	"context"
	"maps"
	"net/http"
	"net/http/httptest"
)

// createTestRequest creates an HTTP request for testing.
func createTestRequest(host string) *http.Request {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+host+testURLPath, nil)
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
		_, err := responseWriter.Write([]byte(testResponseOK))
		if err != nil {
			return
		}
	})
}

// headerCaptureHandler creates a handler that captures the request headers.
func headerCaptureHandler(capturedHeadersRef *map[string][]string) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		*capturedHeadersRef = make(map[string][]string)
		maps.Copy((*capturedHeadersRef), request.Header)
		responseWriter.WriteHeader(http.StatusOK)
		_, err := responseWriter.Write([]byte(testResponseOK))
		if err != nil {
			return
		}
	})
}
