package traefik_conditional_headers

import (
	"context"
	"net/http"
	"testing"
)

// TestMatchesHost tests the host matching function with various scenarios.
func TestMatchesHost(t *testing.T) {
	tests := []struct {
		name         string
		incomingHost string
		ruleHost     string
		expected     bool
	}{
		// Exact matches
		{
			name:         "Exact match simple domain",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostExampleCom,
			expected:     true,
		},
		{
			name:         "Exact match subdomain",
			incomingHost: testHostAPIExampleCom,
			ruleHost:     testHostAPIExampleCom,
			expected:     true,
		},
		{
			name:         "No match different domain",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostOtherCom,
			expected:     false,
		},

		// Port handling
		{
			name:         "Host with port matches rule without port",
			incomingHost: testHostExampleComWithPort,
			ruleHost:     testHostExampleCom,
			expected:     true,
		},

		// Wildcard subdomain matching
		{
			name:         "Wildcard matches subdomain",
			incomingHost: testHostAPIExampleCom,
			ruleHost:     testHostWildcardExampleCom,
			expected:     true,
		},
		{
			name:         "Wildcard matches nested subdomain",
			incomingHost: testHostV1APIExampleCom,
			ruleHost:     testHostWildcardExampleCom,
			expected:     true,
		},
		{
			name:         "Wildcard matches base domain",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostWildcardExampleCom,
			expected:     true,
		},
		{
			name:         "Wildcard no match different domain",
			incomingHost: "api.other.com",
			ruleHost:     testHostWildcardExampleCom,
			expected:     false,
		},

		// Partial matching
		{
			name:         "Contains match partial string",
			incomingHost: testHostMyAPIExampleCom,
			ruleHost:     testHostAPISubstring,
			expected:     true,
		},
		{
			name:         "Contains match at end",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostCOMSubstring,
			expected:     true,
		},
		{
			name:         "Contains match at start",
			incomingHost: testHostAPIExampleCom,
			ruleHost:     testHostAPISubstring,
			expected:     true,
		},
		{
			name:         "Contains no match",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostTestSubstring,
			expected:     false,
		},

		// Edge cases
		{
			name:         "Empty incoming host",
			incomingHost: testHostEmpty,
			ruleHost:     testHostExampleCom,
			expected:     false,
		},
		{
			name:         "Empty rule host",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostEmpty,
			expected:     true, // Empty string contains everything
		},
		{
			name:         "Both empty",
			incomingHost: testHostEmpty,
			ruleHost:     testHostEmpty,
			expected:     true, // Empty string contains empty string
		},
		{
			name:         "Incoming host with port only",
			incomingHost: testHostPortOnly,
			ruleHost:     testHostExampleCom,
			expected:     false,
		},
		{
			name:         "Wildcard rule only",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostWildcardOnly,
			expected:     false,
		},
		{
			name:         "Rule starts with dot",
			incomingHost: testHostExampleCom,
			ruleHost:     testHostDotExampleCom,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesHost(tt.incomingHost, tt.ruleHost)
			if result != tt.expected {
				t.Errorf("matchesHost(%q, %q) = %v, want %v",
					tt.incomingHost, tt.ruleHost, result, tt.expected)
			}
		})
	}
}

// TestCreateConfig tests the configuration creation function.
func TestCreateConfig(t *testing.T) {
	config := CreateConfig()

	if config == nil {
		t.Fatal("CreateConfig() returned nil")
	}

	if config.Rules == nil {
		t.Fatal("CreateConfig() returned config with nil Rules")
	}

	if len(config.Rules) != 0 {
		t.Errorf("CreateConfig() returned config with %d rules, want 0", len(config.Rules))
	}
}

// validateHandlerCreation validates the result of New() function calls.
func validateHandlerCreation(t *testing.T, handler http.Handler, err error, expectError bool) {
	if expectError {
		if err == nil {
			t.Error("Expected error but got none")
		}
		return
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if handler == nil {
		t.Error("New() returned nil handler")
		return
	}

	// Verify the handler implements http.Handler
	if _, ok := handler.(*conditionalHeaders); !ok {
		t.Error("New() did not return conditionalHeaders handler")
	}
}

// TestNew tests the plugin constructor.
func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name:        "Valid config with empty rules",
			config:      &Config{Rules: []Rule{}},
			expectError: false,
		},
		{
			name: "Valid config with rules",
			config: &Config{
				Rules: []Rule{
					{
						Hosts:   []string{testHostExampleCom},
						Headers: map[string]string{testHeaderXTest: testValueValue},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "Nil config",
			config:      &Config{Rules: []Rule{}},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := mockNextHandler()
			handler, err := New(context.Background(), next, tt.config, testPluginNameWithDash)
			validateHandlerCreation(t, handler, err, tt.expectError)
		})
	}
}

// TestConditionalHeadersFields tests the internal state of the handler.
func TestConditionalHeadersFields(t *testing.T) {
	rules := []Rule{
		{
			Hosts:   []string{testHostExampleCom},
			Headers: map[string]string{testHeaderXTest: testValueValue},
		},
	}

	config := &Config{Rules: rules}
	next := mockNextHandler()
	handler, err := New(context.Background(), next, config, testPluginNameWithDash)

	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	conditionalHandler, ok := handler.(*conditionalHeaders)
	if !ok {
		t.Fatalf("Expected handler to be of type *conditionalHeaders, got %T", handler)
	}

	if conditionalHandler.next == nil {
		t.Error("Handler next field is nil")
	}

	if len(conditionalHandler.rules) != 1 {
		t.Errorf("Handler rules length: got %d, want 1", len(conditionalHandler.rules))
	}

	if conditionalHandler.name != testPluginNameWithDash {
		t.Errorf("Handler name: got %q, want %q", conditionalHandler.name, testPluginNameWithDash)
	}

	// Verify rule content
	if len(conditionalHandler.rules[0].Hosts) != 1 {
		t.Errorf("Rule hosts length: got %d, want 1", len(conditionalHandler.rules[0].Hosts))
	}

	if conditionalHandler.rules[0].Hosts[0] != testHostExampleCom {
		t.Errorf("Rule host: got %q, want %q", conditionalHandler.rules[0].Hosts[0], testHostExampleCom)
	}

	if conditionalHandler.rules[0].Headers[testHeaderXTest] != testValueValue {
		t.Errorf("Rule header: got %q, want %q", conditionalHandler.rules[0].Headers[testHeaderXTest], testValueValue)
	}
}

// BenchmarkMatchesHost benchmarks the host matching function.
func BenchmarkMatchesHost(b *testing.B) {
	testCases := []struct {
		name         string
		incomingHost string
		ruleHost     string
	}{
		{"Exact match", testHostExampleCom, testHostExampleCom},
		{"Wildcard match", testHostAPIExampleCom, testHostWildcardExampleCom},
		{"Contains match", testHostAPIExampleCom, testHostAPISubstring},
		{"No match", testHostExampleCom, testHostOtherCom},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matchesHost(tc.incomingHost, tc.ruleHost)
			}
		})
	}
}
