package traefik_conditional_headers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// verifyExpectedHeaders checks that the expected headers were set correctly.
func verifyExpectedHeaders(t *testing.T, capturedHeaders map[string][]string, expectedHeaders map[string]string) {
	if len(capturedHeaders) == 0 {
		t.Error("Expected headers to be set, but none were captured")
	}

	for expectedKey, expectedValue := range expectedHeaders {
		values, exists := capturedHeaders[expectedKey]
		if !exists {
			t.Errorf("Expected header %q to be set", expectedKey)
			continue
		}
		if len(values) != 1 || values[0] != expectedValue {
			t.Errorf("Header %q: expected %q, got %v", expectedKey, expectedValue, values)
		}
	}
}

// verifyNoUnexpectedHeaders checks that no unexpected headers were set.
func verifyNoUnexpectedHeaders(t *testing.T, capturedHeaders map[string][]string) {
	for header := range capturedHeaders {
		// Skip standard headers that might be set by the test infrastructure
		if header != "User-Agent" && header != "Accept-Encoding" {
			t.Errorf("Unexpected header %q was set", header)
		}
	}
}

// executeTestRequest creates a handler and executes a test request.
func executeTestRequest(t *testing.T, rules []Rule, requestHost string) (map[string][]string, *httptest.ResponseRecorder) {
	var capturedHeaders map[string][]string
	next := headerCaptureHandler(&capturedHeaders)

	config := &Config{Rules: rules}
	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest(requestHost)
	responseRecorder := createTestResponse()

	handler.ServeHTTP(responseRecorder, req)

	// Verify response
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	return capturedHeaders, responseRecorder
}

// Test_conditionalHeaders_ServeHTTP tests the main HTTP handler functionality.
func Test_conditionalHeaders_ServeHTTP(t *testing.T) {
	tests := []struct {
		name            string
		rules           []Rule
		requestHost     string
		expectedHeaders map[string]string
		shouldMatch     bool
	}{
		{
			name:            "No matching rules",
			rules:           []Rule{{Hosts: []string{"example.com"}, Headers: map[string]string{"X-Test": "value"}}},
			requestHost:     "other.com",
			expectedHeaders: map[string]string{},
			shouldMatch:     false,
		},
		{
			name:            "Exact host match",
			rules:           []Rule{{Hosts: []string{"example.com"}, Headers: map[string]string{"X-Custom": "test-value"}}},
			requestHost:     "example.com",
			expectedHeaders: map[string]string{"X-Custom": "test-value"},
			shouldMatch:     true,
		},
		{
			name:            "Wildcard host match",
			rules:           []Rule{{Hosts: []string{"*.example.com"}, Headers: map[string]string{"X-Environment": "development"}}},
			requestHost:     "api.example.com",
			expectedHeaders: map[string]string{"X-Environment": "development"},
			shouldMatch:     true,
		},
		{
			name: "Multiple hosts in rule match",
			rules: []Rule{
				{
					Hosts:   []string{"example.com", "api.example.com", "test.com"},
					Headers: map[string]string{"X-Service": "api", "X-Version": "v1"},
				},
			},
			requestHost:     "api.example.com",
			expectedHeaders: map[string]string{"X-Service": "api", "X-Version": "v1"},
			shouldMatch:     true,
		},
		{
			name: "First rule wins - multiple matching rules",
			rules: []Rule{
				{Hosts: []string{"*.example.com"}, Headers: map[string]string{"X-First": "should-win"}},
				{Hosts: []string{"api.example.com"}, Headers: map[string]string{"X-Second": "should-not-win"}},
			},
			requestHost:     "api.example.com",
			expectedHeaders: map[string]string{"X-First": "should-win"},
			shouldMatch:     true,
		},
		{
			name:            "Host with port matches",
			rules:           []Rule{{Hosts: []string{"example.com"}, Headers: map[string]string{"X-Port-Test": "success"}}},
			requestHost:     "example.com:8080",
			expectedHeaders: map[string]string{"X-Port-Test": "success"},
			shouldMatch:     true,
		},
		{
			name:            "Partial match works",
			rules:           []Rule{{Hosts: []string{"api"}, Headers: map[string]string{"X-Partial": "match"}}},
			requestHost:     "my-api.example.com",
			expectedHeaders: map[string]string{"X-Partial": "match"},
			shouldMatch:     true,
		},
		{
			name:            "Empty rules list",
			rules:           []Rule{},
			requestHost:     "example.com",
			expectedHeaders: map[string]string{},
			shouldMatch:     false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			capturedHeaders, _ := executeTestRequest(t, testCase.rules, testCase.requestHost)

			if testCase.shouldMatch {
				verifyExpectedHeaders(t, capturedHeaders, testCase.expectedHeaders)
			} else {
				verifyNoUnexpectedHeaders(t, capturedHeaders)
			}
		})
	}
}

// Test_conditionalHeaders_ServeHTTP_requestFlow tests that requests flow properly to the next handler.
func Test_conditionalHeaders_ServeHTTP_requestFlow(t *testing.T) {
	// Track if the next handler was called
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		nextCalled = true

		// Set a response to verify the handler was called correctly
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte("Next handler called"))
		if err != nil {
			return
		}
	})

	config := &Config{
		Rules: []Rule{
			{
				Hosts:   []string{"example.com"},
				Headers: map[string]string{"X-Test": "value"},
			},
		},
	}

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest("example.com")
	responseRecorder := createTestResponse()

	handler.ServeHTTP(responseRecorder, req)

	// Verify the next handler was called
	if !nextCalled {
		t.Error("Next handler was not called")
	}

	// Verify response
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	expectedBody := "Next handler called"
	if responseRecorder.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, responseRecorder.Body.String())
	}

	// Verify headers were set
	if req.Header.Get("X-Test") != "value" {
		t.Errorf("Expected header X-Test to be 'value', got %q", req.Header.Get("X-Test"))
	}
}

// Test_conditionalHeaders_ServeHTTP_noRules tests behavior when no rules are configured.
func Test_conditionalHeaders_ServeHTTP_noRules(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		nextCalled = true
		rw.WriteHeader(http.StatusOK)
	})

	config := &Config{Rules: []Rule{}}
	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest("example.com")
	responseRecorder := createTestResponse()

	handler.ServeHTTP(responseRecorder, req)

	if !nextCalled {
		t.Error("Next handler was not called")
	}

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

// Test_conditionalHeaders_ServeHTTP_nilNext tests behavior when next handler is nil.
func Test_conditionalHeaders_ServeHTTP_nilNext(t *testing.T) {
	config := &Config{
		Rules: []Rule{
			{
				Hosts:   []string{"example.com"},
				Headers: map[string]string{"X-Test": "value"},
			},
		},
	}

	handler, err := New(context.Background(), nil, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest("example.com")
	responseRecorder := createTestResponse()

	// This should NOT panic - nil next handlers are handled gracefully for yaegi compatibility
	handler.ServeHTTP(responseRecorder, req)

	// Verify that headers were set despite next being nil
	if req.Header.Get("X-Test") != "value" {
		t.Errorf("Expected header X-Test to be 'value', got %q", req.Header.Get("X-Test"))
	}
}

// Test_conditionalHeaders_multipleHeaders tests multiple headers being set.
func Test_conditionalHeaders_multipleHeaders(t *testing.T) {
	headers := map[string]string{
		"X-Service":     "api",
		"X-Version":     "v2",
		"X-Environment": "production",
		"X-Custom":      "test-value",
	}

	config := &Config{
		Rules: []Rule{
			{
				Hosts:   []string{"api.example.com"},
				Headers: headers,
			},
		},
	}

	next := mockNextHandler()
	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest("api.example.com")
	responseRecorder := createTestResponse()

	handler.ServeHTTP(responseRecorder, req)

	// Verify all headers were set correctly
	for key, expectedValue := range headers {
		actualValue := req.Header.Get(key)
		if actualValue != expectedValue {
			t.Errorf("Header %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

// Benchmark_conditionalHeaders_ServeHTTP benchmarks the ServeHTTP method.
func Benchmark_conditionalHeaders_ServeHTTP(b *testing.B) {
	config := &Config{
		Rules: []Rule{
			{
				Hosts:   []string{"example.com", "*.example.com"},
				Headers: map[string]string{"X-Test": "value"},
			},
		},
	}

	next := mockNextHandler()
	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		b.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest("api.example.com")
	responseRecorder := createTestResponse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		responseRecorder.Body.Reset()
		handler.ServeHTTP(responseRecorder, req)
	}
}
