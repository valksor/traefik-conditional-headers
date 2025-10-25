package traefik_conditional_headers

import (
	"context"
	"testing"
)

// Test_matchesHost tests the host matching function with various scenarios.
func Test_matchesHost(t *testing.T) {
	tests := []struct {
		name         string
		incomingHost string
		ruleHost     string
		expected     bool
	}{
		// Exact matches
		{
			name:         "Exact match simple domain",
			incomingHost: "example.com",
			ruleHost:     "example.com",
			expected:     true,
		},
		{
			name:         "Exact match subdomain",
			incomingHost: "api.example.com",
			ruleHost:     "api.example.com",
			expected:     true,
		},
		{
			name:         "No match different domain",
			incomingHost: "example.com",
			ruleHost:     "other.com",
			expected:     false,
		},

		// Port handling
		{
			name:         "Host with port matches rule without port",
			incomingHost: "example.com:8080",
			ruleHost:     "example.com",
			expected:     true,
		},
	
		// Wildcard subdomain matching
		{
			name:         "Wildcard matches subdomain",
			incomingHost: "api.example.com",
			ruleHost:     "*.example.com",
			expected:     true,
		},
		{
			name:         "Wildcard matches nested subdomain",
			incomingHost: "v1.api.example.com",
			ruleHost:     "*.example.com",
			expected:     true,
		},
		{
			name:         "Wildcard matches base domain",
			incomingHost: "example.com",
			ruleHost:     "*.example.com",
			expected:     true,
		},
		{
			name:         "Wildcard no match different domain",
			incomingHost: "api.other.com",
			ruleHost:     "*.example.com",
			expected:     false,
		},

		// Partial matching
		{
			name:         "Contains match partial string",
			incomingHost: "my-api.example.com",
			ruleHost:     "api",
			expected:     true,
		},
		{
			name:         "Contains match at end",
			incomingHost: "example.com",
			ruleHost:     "com",
			expected:     true,
		},
		{
			name:         "Contains match at start",
			incomingHost: "api.example.com",
			ruleHost:     "api",
			expected:     true,
		},
		{
			name:         "Contains no match",
			incomingHost: "example.com",
			ruleHost:     "test",
			expected:     false,
		},

		// Edge cases
		{
			name:         "Empty incoming host",
			incomingHost: "",
			ruleHost:     "example.com",
			expected:     false,
		},
		{
			name:         "Empty rule host",
			incomingHost: "example.com",
			ruleHost:     "",
			expected:     true, // Empty string contains everything
		},
		{
			name:         "Both empty",
			incomingHost: "",
			ruleHost:     "",
			expected:     true, // Empty string contains empty string
		},
		{
			name:         "Incoming host with port only",
			incomingHost: ":8080",
			ruleHost:     "example.com",
			expected:     false,
		},
		{
			name:         "Wildcard rule only",
			incomingHost: "example.com",
			ruleHost:     "*",
			expected:     false,
		},
		{
			name:         "Rule starts with dot",
			incomingHost: "example.com",
			ruleHost:     ".example.com",
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
			name:   "Valid config with rules",
			config: &Config{
				Rules: []Rule{
					{
						Hosts:   []string{"example.com"},
						Headers: map[string]string{"X-Test": "value"},
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
			handler, err := New(context.Background(), next, tt.config, "test-plugin")

			if tt.expectError {
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
		})
	}
}

// Test_conditionalHeaders_fields tests the internal state of the handler.
func Test_conditionalHeaders_fields(t *testing.T) {
	rules := []Rule{
		{
			Hosts:   []string{"example.com"},
			Headers: map[string]string{"X-Test": "value"},
		},
	}

	config := &Config{Rules: rules}
	next := mockNextHandler()
	handler, err := New(context.Background(), next, config, "test-plugin")

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

	if conditionalHandler.name != "test-plugin" {
		t.Errorf("Handler name: got %q, want %q", conditionalHandler.name, "test-plugin")
	}

	// Verify rule content
	if len(conditionalHandler.rules[0].Hosts) != 1 {
		t.Errorf("Rule hosts length: got %d, want 1", len(conditionalHandler.rules[0].Hosts))
	}

	if conditionalHandler.rules[0].Hosts[0] != "example.com" {
		t.Errorf("Rule host: got %q, want %q", conditionalHandler.rules[0].Hosts[0], "example.com")
	}

	if conditionalHandler.rules[0].Headers["X-Test"] != "value" {
		t.Errorf("Rule header: got %q, want %q", conditionalHandler.rules[0].Headers["X-Test"], "value")
	}
}

// Benchmark_matchesHost benchmarks the host matching function.
func Benchmark_matchesHost(b *testing.B) {
	testCases := []struct {
		name         string
		incomingHost string
		ruleHost     string
	}{
		{"Exact match", "example.com", "example.com"},
		{"Wildcard match", "api.example.com", "*.example.com"},
		{"Contains match", "api.example.com", "api"},
		{"No match", "example.com", "other.com"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matchesHost(tc.incomingHost, tc.ruleHost)
			}
		})
	}
}