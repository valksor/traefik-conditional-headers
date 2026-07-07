package traefik_conditional_headers

import (
	"context"
	"net/http"
	"testing"
)

// assertExpectedHeaders validates that expected headers are present in actual headers.
func assertExpectedHeaders(t *testing.T, host string, expectedHeaders map[string]string, actualHeaders map[string][]string) {
	t.Helper()

	for expectedKey, expectedValue := range expectedHeaders {
		actualValue, exists := actualHeaders[expectedKey]
		if !exists {
			t.Errorf("Expected header %q to be set for host %q", expectedKey, host)

			continue
		}

		// For simplicity, we'll check the first value (headers are stored as slices)
		if len(actualValue) == 0 || actualValue[0] != expectedValue {
			t.Errorf("Header %q for host %q: expected %q, got %q",
				expectedKey, host, expectedValue, actualValue)
		}
	}
}

// validateTestRequest executes a test request and validates expected headers.
func validateTestRequest(t *testing.T, config *Config, testReq testRequest) {
	t.Helper()

	headers := makeTestRequest(t, config, testReq.host)
	assertExpectedHeaders(t, testReq.host, testReq.expectedHeaders, headers)
}

// runTestCase executes all test requests for a given test case.
func runTestCase(t *testing.T, testCase struct {
	name         string
	config       *Config
	testRequests []testRequest
	description  string
},
) {
	t.Helper()

	t.Logf("Testing scenario: %s", testCase.description)

	for _, testReq := range testCase.testRequests {
		t.Run(testReq.description, func(t *testing.T) {
			validateTestRequest(t, testCase.config, testReq)
		})
	}
}

// TestIntegrationRealWorldScenarios tests realistic configuration scenarios.
func TestIntegrationRealWorldScenarios(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		testRequests []testRequest
		description  string
	}{
		{
			name:        "Multi-environment setup",
			description: "Test configuration for different environments with subdomains",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{testHostWildcardDevExample, testHostAPIDevExampleCom, testHostWildcardStagingExample, testHostStagingAPIExample},
						Headers: map[string]string{
							testHeaderXEnvironment: testValueDevelopment,
							testHeaderXDebug:       testValueTrue,
							testHeaderXCORS:        "*",
						},
					},
					{
						Hosts: []string{testHostAPIExampleCom, testHostProductionAPIExample},
						Headers: map[string]string{
							testHeaderXEnvironment:  testValueProduction,
							testHeaderXCacheControl: testValueNoCache,
						},
					},
					{
						Hosts: []string{testHostWildcardExampleCom},
						Headers: map[string]string{
							testHeaderXService: testValueAPIGateway,
							testHeaderXVersion: testValueV210,
						},
					},
				},
			},
			testRequests: []testRequest{
				{
					host: testHostAPIDevExampleCom,
					expectedHeaders: map[string]string{
						testHeaderXEnvironment: testValueDevelopment,
						testHeaderXDebug:       testValueTrue,
						"X-Cors":               "*",
					},
					description: "Development API gets debug headers",
				},
				{
					host: testHostAPIExampleCom,
					expectedHeaders: map[string]string{
						testHeaderXEnvironment:  testValueProduction,
						testHeaderXCacheControl: testValueNoCache,
					},
					description: "Production API gets cache control headers",
				},
				{
					host: testHostTestExampleCom,
					expectedHeaders: map[string]string{
						testHeaderXService: testValueAPIGateway,
						testHeaderXVersion: testValueV210,
					},
					description: "Test subdomain gets wildcard rule headers",
				},
				{
					host:            testHostOtherCom,
					expectedHeaders: map[string]string{},
					description:     "External domain gets no headers",
				},
			},
		},
		{
			name:        "Microservice routing with authentication",
			description: "Test microservice setup with authentication and versioning",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{testHostAuthServiceLocal, testHostLoginServiceLocal},
						Headers: map[string]string{
							testHeaderXServiceType:  testValueAuthentication,
							testHeaderXAuthRequired: testValueFalse,
							testHeaderXPublic:       testValueTrue,
						},
					},
					{
						Hosts: []string{testHostUsersAPIServiceLocal, testHostProfileAPIService},
						Headers: map[string]string{
							testHeaderXServiceType:  testValueUserManagement,
							testHeaderXAuthRequired: testValueTrue,
							testHeaderXRateLimit:    testValue1000PerHour,
						},
					},
					{
						Hosts: []string{testHostAdminServiceLocal},
						Headers: map[string]string{
							testHeaderXServiceType:  testValueAdministration,
							testHeaderXAuthRequired: testValueTrue,
							testHeaderXRoleRequired: testValueAdmin,
							testHeaderXAuditLog:     testValueTrue,
						},
					},
					{
						Hosts: []string{testHostWildcardServiceLocal},
						Headers: map[string]string{
							testHeaderXCluster:    testValueProductionCluster,
							testHeaderXDatacenter: testValueUSWest2,
						},
					},
				},
			},
			testRequests: []testRequest{
				{
					host: testHostAuthServiceLocal,
					expectedHeaders: map[string]string{
						testHeaderXServiceType:  testValueAuthentication,
						testHeaderXAuthRequired: testValueFalse,
						testHeaderXPublic:       testValueTrue,
					},
					description: "Auth service gets public headers",
				},
				{
					host: testHostUsersAPIServiceLocal,
					expectedHeaders: map[string]string{
						testHeaderXServiceType:  testValueUserManagement,
						testHeaderXAuthRequired: testValueTrue,
						testHeaderXRateLimit:    testValue1000PerHour,
					},
					description: "Users API gets authentication and rate limiting headers",
				},
				{
					host: testHostAdminServiceLocal,
					expectedHeaders: map[string]string{
						testHeaderXServiceType:  testValueAdministration,
						testHeaderXAuthRequired: testValueTrue,
						testHeaderXRoleRequired: testValueAdmin,
						testHeaderXAuditLog:     testValueTrue,
					},
					description: "Admin service gets strict security headers",
				},
				{
					host: testHostCacheServiceLocal,
					expectedHeaders: map[string]string{
						testHeaderXCluster:    testValueProductionCluster,
						testHeaderXDatacenter: testValueUSWest2,
					},
					description: "Generic service gets cluster headers only",
				},
			},
		},
		{
			name:        "CDN and static file hosting",
			description: "Test CDN setup with different file types and caching strategies",
			config: &Config{
				Rules: []Rule{
					{
						Hosts: []string{testHostImagesCDNExample, testHostImgCDNExample},
						Headers: map[string]string{
							testHeaderXContentType:  testValueImage,
							testHeaderXCacheControl: testValuePublicImmutable,
							testHeaderXCompress:     testValueGzip,
						},
					},
					{
						Hosts: []string{testHostAPIExampleCom},
						Headers: map[string]string{
							testHeaderXContentType:  testValueAPI2,
							testHeaderXCacheControl: testValueNoCache,
							testHeaderXRateLimit:    testValue10000PerHour,
						},
					},
					{
						Hosts: []string{testHostCDNExample, testHostStaticExample, testHostAssetsExample},
						Headers: map[string]string{
							testHeaderXContentType:  testValueStatic,
							testHeaderXCacheControl: testValuePublicMaxAge,
							testHeaderXCDN:          testValueEnabled,
						},
					},
				},
			},
			testRequests: []testRequest{
				{
					host: testHostCDNExample,
					expectedHeaders: map[string]string{
						testHeaderXContentType:  testValueStatic,
						testHeaderXCacheControl: testValuePublicMaxAge,
						"X-Cdn":                 testValueEnabled,
					},
					description: "Main CDN gets static content headers",
				},
				{
					host: testHostImagesCDNExample,
					expectedHeaders: map[string]string{
						testHeaderXContentType:  testValueImage,
						testHeaderXCacheControl: testValuePublicImmutable,
						testHeaderXCompress:     testValueGzip,
					},
					description: "Image CDN gets aggressive caching headers",
				},
				{
					host: testHostAPIExampleCom,
					expectedHeaders: map[string]string{
						testHeaderXContentType:  testValueAPI2,
						testHeaderXCacheControl: testValueNoCache,
						testHeaderXRateLimit:    testValue10000PerHour,
					},
					description: "API gets no-cache headers",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			runTestCase(t, testCase)
		})
	}
}

// TestIntegrationEdgeCases tests edge cases and boundary conditions.
func TestIntegrationEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		host        string
		description string
		expectPanic bool
	}{
		{
			name:        "Overlapping wildcard rules",
			description: "Test behavior with overlapping wildcard patterns",
			config: &Config{
				Rules: []Rule{
					{
						Hosts:   []string{testHostWildcardExampleCom},
						Headers: map[string]string{testHeaderXLevel: testValueSubdomain},
					},
					{
						Hosts:   []string{testHostAPIExampleComWildcard},
						Headers: map[string]string{testHeaderXLevel: testValueAPISubdomain},
					},
				},
			},
			host:        testHostV1APIExampleCom,
			expectPanic: false,
		},
		{
			name:        "Unicode and special characters",
			description: "Test hosts with unicode characters and special cases",
			config: &Config{
				Rules: []Rule{
					{
						Hosts:   []string{testHostTestAccent},
						Headers: map[string]string{testHeaderXUnicode: testValueTest},
					},
					{
						Hosts:   []string{testHostWildcardExampleCom},
						Headers: map[string]string{testHeaderXWildcard: testValueMatch},
					},
				},
			},
			host:        testHostTestAccent,
			expectPanic: false,
		},
		{
			name:        "Empty configuration",
			description: "Test with completely empty configuration",
			config:      &Config{Rules: []Rule{}},
			host:        testHostExampleCom,
			expectPanic: false,
		},
		{
			name:        "Configuration with empty rule",
			description: "Test with rule that has empty hosts and headers",
			config: &Config{
				Rules: []Rule{
					{
						Hosts:   []string{},
						Headers: map[string]string{},
					},
					{
						Hosts:   []string{testHostExampleCom},
						Headers: map[string]string{testHeaderXTest: testValueValue},
					},
				},
			},
			host:        testHostExampleCom,
			expectPanic: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Logf("Testing edge case: %s", testCase.description)

			if testCase.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but none occurred")
					}
				}()
			}

			// This should not panic for normal cases
			headers := makeTestRequest(t, testCase.config, testCase.host)

			// Just verify we got some result (could be empty map)
			if headers == nil {
				t.Error("Expected headers map, got nil")
			}
		})
	}
}

// TestIntegrationPerformanceComplexConfig tests performance with complex configurations.
func TestIntegrationPerformanceComplexConfig(t *testing.T) {
	// Create a complex configuration with many rules
	hostSuffixes := []string{".com", ".net", ".org", ".io", ".dev"}
	subdomains := []string{"api", "admin", "cdn", "static", "auth", "users", "blog", "shop"}
	rules := make([]Rule, 0, len(subdomains)*len(hostSuffixes)+1)

	// Generate many rules
	for suffixIndex, suffix := range hostSuffixes {
		for subdomainIndex, subdomain := range subdomains {
			domain := "example" + suffix
			rules = append(rules, Rule{
				Hosts: []string{subdomain + "." + domain},
				Headers: map[string]string{
					"X-Service": subdomain,
					"X-Domain":  domain,
					"X-Rule-ID": string(rune(suffixIndex*10 + subdomainIndex)),
				},
			})
		}
	}

	// Add some wildcard rules
	rules = append(rules, Rule{
		Hosts: []string{testHostWildcardExampleCom, testHostWildcardExampleNet},
		Headers: map[string]string{
			testHeaderXWildcard: testValuePrimary,
		},
	})

	config := &Config{Rules: rules}

	t.Logf("Testing with %d rules", len(rules))

	// Test various hosts to ensure performance is acceptable
	testHosts := []string{
		testHostAPIExampleCom,
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
				if start > 2000 { // 2000 bytes is a reasonable threshold for complex rule matching
					t.Errorf("Too many allocations for host %s: %.0f bytes", host, start)
				}
			}
		})
	}
}

// Helper types and functions for integration tests

type testRequest struct {
	host            string
	expectedHeaders map[string]string
	description     string
}

// makeTestRequest creates and executes a test request, returning the headers.
func makeTestRequest(t *testing.T, config *Config, host string) map[string][]string {
	t.Helper()

	var capturedHeaders map[string][]string
	next := headerCaptureHandler(&capturedHeaders)

	handler, err := New(context.Background(), next, config, testPluginNameIntegration)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	req := createTestRequest(host)
	rw := createTestResponse()

	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}

	return capturedHeaders
}
