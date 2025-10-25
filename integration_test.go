package traefik_conditional_headers

import (
	"context"
	"net/http"
	"testing"
)

// TestIntegration_RealWorldScenarios tests realistic configuration scenarios
func TestIntegration_RealWorldScenarios(t *testing.T) {
	tests := []struct {
		name            string
		config          *Config
		testRequests    []testRequest
		description     string
	}{
		{
			name: "Multi-environment setup",
			description: "Test configuration for different environments with subdomains",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{"*.dev.example.com", "api.dev.example.com", "*.staging.example.com", "staging-api.example.com"},
						Headers: map[string]string{
							"X-Environment": "development",
							"X-Debug":       "true",
							"X-CORS":        "*",
						},
					},
					{
						Hosts: []string{"api.example.com", "production-api.example.com"},
						Headers: map[string]string{
							"X-Environment": "production",
							"X-Cache-Control": "no-cache",
						},
					},
					{
						Hosts: []string{"*.example.com"},
						Headers: map[string]string{
							"X-Service": "api-gateway",
							"X-Version": "v2.1.0",
						},
					},
				},
			},
			testRequests: []testRequest{
				{
					host: "api.dev.example.com",
					expectedHeaders: map[string]string{
						"X-Environment": "development",
						"X-Debug":       "true",
						"X-Cors":        "*",
					},
					description: "Development API gets debug headers",
				},
				{
					host: "api.example.com",
					expectedHeaders: map[string]string{
						"X-Environment":  "production",
						"X-Cache-Control": "no-cache",
					},
					description: "Production API gets cache control headers",
				},
				{
					host: "test.example.com",
					expectedHeaders: map[string]string{
						"X-Service": "api-gateway",
						"X-Version": "v2.1.0",
					},
					description: "Test subdomain gets wildcard rule headers",
				},
				{
					host: "other.com",
					expectedHeaders: map[string]string{},
					description: "External domain gets no headers",
				},
			},
		},
		{
			name: "Microservice routing with authentication",
			description: "Test microservice setup with authentication and versioning",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{"auth.service.local", "login.service.local"},
						Headers: map[string]string{
							"X-Service-Type": "authentication",
							"X-Auth-Required": "false",
							"X-Public":       "true",
						},
					},
					{
						Hosts: []string{"users.api.service.local", "profile.api.service.local"},
						Headers: map[string]string{
							"X-Service-Type": "user-management",
							"X-Auth-Required": "true",
							"X-Rate-Limit":    "1000/hour",
						},
					},
					{
						Hosts: []string{"admin.service.local"},
						Headers: map[string]string{
							"X-Service-Type": "administration",
							"X-Auth-Required": "true",
							"X-Role-Required": "admin",
							"X-Audit-Log":     "true",
						},
					},
					{
						Hosts: []string{"*.service.local"},
						Headers: map[string]string{
							"X-Cluster": "production",
							"X-Datacenter": "us-west-2",
						},
					},
				},
			},
			testRequests: []testRequest{
				{
					host: "auth.service.local",
					expectedHeaders: map[string]string{
						"X-Service-Type": "authentication",
						"X-Auth-Required": "false",
						"X-Public":       "true",
					},
					description: "Auth service gets public headers",
				},
				{
					host: "users.api.service.local",
					expectedHeaders: map[string]string{
						"X-Service-Type": "user-management",
						"X-Auth-Required": "true",
						"X-Rate-Limit":    "1000/hour",
					},
					description: "Users API gets authentication and rate limiting headers",
				},
				{
					host: "admin.service.local",
					expectedHeaders: map[string]string{
						"X-Service-Type": "administration",
						"X-Auth-Required": "true",
						"X-Role-Required": "admin",
						"X-Audit-Log":     "true",
					},
					description: "Admin service gets strict security headers",
				},
				{
					host: "cache.service.local",
					expectedHeaders: map[string]string{
						"X-Cluster": "production",
						"X-Datacenter": "us-west-2",
					},
					description: "Generic service gets cluster headers only",
				},
			},
		},
		{
			name: "CDN and static file hosting",
			description: "Test CDN setup with different file types and caching strategies",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{"images.cdn.example.com", "img.cdn.example.com"},
						Headers: map[string]string{
							"X-Content-Type": "image",
							"X-Cache-Control": "public, max-age=31536000, immutable",
							"X-Compress": "gzip",
						},
					},
					{
						Hosts: []string{"api.example.com"},
						Headers: map[string]string{
							"X-Content-Type": "api",
							"X-Cache-Control": "no-cache",
							"X-Rate-Limit": "10000/hour",
						},
					},
					{
						Hosts: []string{"cdn.example.com", "static.example.com", "assets.example.com"},
						Headers: map[string]string{
							"X-Content-Type": "static",
							"X-Cache-Control": "public, max-age=31536000",
							"X-CDN": "enabled",
						},
					},
			},
		},
		testRequests: []testRequest{
				{
					host: "cdn.example.com",
					expectedHeaders: map[string]string{
						"X-Content-Type": "static",
						"X-Cache-Control": "public, max-age=31536000",
						"X-Cdn": "enabled",
					},
					description: "Main CDN gets static content headers",
				},
				{
					host: "images.cdn.example.com",
					expectedHeaders: map[string]string{
						"X-Content-Type": "image",
						"X-Cache-Control": "public, max-age=31536000, immutable",
						"X-Compress": "gzip",
					},
					description: "Image CDN gets aggressive caching headers",
				},
				{
					host: "api.example.com",
					expectedHeaders: map[string]string{
						"X-Content-Type": "api",
						"X-Cache-Control": "no-cache",
						"X-Rate-Limit": "10000/hour",
					},
					description: "API gets no-cache headers",
				},
		},
	},
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing scenario: %s", tt.description)

			for _, testReq := range tt.testRequests {
				t.Run(testReq.description, func(t *testing.T) {
					// Test the request
					headers := makeTestRequest(t, tt.config, testReq.host)

					// Verify expected headers are set
					for expectedKey, expectedValue := range testReq.expectedHeaders {
						actualValue, exists := headers[expectedKey]
						if !exists {
							t.Errorf("Expected header %q to be set for host %q", expectedKey, testReq.host)
							continue
						}

						// For simplicity, we'll check the first value (headers are stored as slices)
						if len(actualValue) == 0 || actualValue[0] != expectedValue {
							t.Errorf("Header %q for host %q: expected %q, got %q",
								expectedKey, testReq.host, expectedValue, actualValue)
						}
					}
				})
			}
		})
	}
}

// TestIntegration_EdgeCases tests edge cases and boundary conditions
func TestIntegration_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		host        string
		description string
		expectPanic bool
	}{
		{
			name: "Overlapping wildcard rules",
			description: "Test behavior with overlapping wildcard patterns",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{"*.example.com"},
						Headers: map[string]string{"X-Level": "subdomain"},
					},
					{
						Hosts: []string{"*.api.example.com"},
						Headers: map[string]string{"X-Level": "api-subdomain"},
					},
				},
			},
			host: "v1.api.example.com",
			expectPanic: false,
		},
		{
			name: "Unicode and special characters",
			description: "Test hosts with unicode characters and special cases",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{"tést.example.com"},
						Headers: map[string]string{"X-Unicode": "test"},
					},
					{
						Hosts: []string{"*.example.com"},
						Headers: map[string]string{"X-Wildcard": "match"},
					},
				},
			},
			host: "tést.example.com",
			expectPanic: false,
		},
		{
			name: "Empty configuration",
			description: "Test with completely empty configuration",
			config: &Config{},
			host: "example.com",
			expectPanic: false,
		},
		{
			name: "Configuration with empty rule",
			description: "Test with rule that has empty hosts and headers",
			config: &Config{
				Rules: []Rule{
					{
						Hosts:   []string{},
						Headers: map[string]string{},
					},
					{
						Hosts: []string{"example.com"},
						Headers: map[string]string{"X-Test": "value"},
					},
				},
			},
			host: "example.com",
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing edge case: %s", tt.description)

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			// This should not panic for normal cases
			headers := makeTestRequest(t, tt.config, tt.host)

			// Just verify we got some result (could be empty map)
			if headers == nil {
				t.Error("Expected headers map, got nil")
			}
		})
	}
}

// TestIntegration_PerformanceComplexConfig tests performance with complex configurations
func TestIntegration_PerformanceComplexConfig(t *testing.T) {
	// Create a complex configuration with many rules
	var rules []Rule
	hostSuffixes := []string{".com", ".net", ".org", ".io", ".dev"}
	subdomains := []string{"api", "admin", "cdn", "static", "auth", "users", "blog", "shop"}

	// Generate many rules
	for i, suffix := range hostSuffixes {
		for j, subdomain := range subdomains {
			domain := "example" + suffix
			rules = append(rules, Rule{
				Hosts: []string{subdomain + "." + domain},
				Headers: map[string]string{
					"X-Service": subdomain,
					"X-Domain":  domain,
					"X-Rule-ID": string(rune(i*10 + j)),
				},
			})
		}
	}

	// Add some wildcard rules
	rules = append(rules, Rule{
		Hosts: []string{"*.example.com", "*.example.net"},
		Headers: map[string]string{
			"X-Wildcard": "primary",
		},
	})

	config := &Config{Rules: rules}

	t.Logf("Testing with %d rules", len(rules))

	// Test various hosts to ensure performance is acceptable
	testHosts := []string{
		"api.example.com",
		"admin.example.org",
		"cdn.example.io",
		"users.example.dev",
		"unknown.example.com",
		"other.domain.com",
	}

	for _, host := range testHosts {
		t.Run("Performance test for "+host, func(t *testing.T) {
			// Run multiple iterations to check performance
			for i := 0; i < 100; i++ {
				start := testing.AllocsPerRun(1, func() {
					makeTestRequest(t, config, host)
				})
				// Ensure we're not allocating too much per request
				if start > 1000 { // 1000 bytes is a reasonable threshold
					t.Errorf("Too many allocations for host %s: %.0f bytes", host, start)
				}
			}
		})
	}
}

// Helper types and functions for integration tests

type testRequest struct {
	host             string
	expectedHeaders  map[string]string
	description      string
}

// makeTestRequest creates and executes a test request, returning the headers
func makeTestRequest(t *testing.T, config *Config, host string) map[string][]string {
	var capturedHeaders map[string][]string
	next := headerCaptureHandler(&capturedHeaders)

	handler, err := New(context.Background(), next, config, "integration-test")
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest(host, "/test")
	rw := createTestResponse()

	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}

	return capturedHeaders
}