// Package traefik_conditional_headers provides a Traefik middleware plugin for setting HTTP headers
// conditionally based on the incoming request hostname.
//
// The plugin supports multiple matching strategies:
//   - Exact host matching (e.g., "api.example.com")
//   - Wildcard subdomain matching (e.g., "*.example.com")
//   - Partial string matching (e.g., "api" matches "my-api.example.com")
//
// Rules are evaluated in order, and the first matching rule wins. Multiple hosts can share
// the same header configuration within a single rule.
//
// Example configuration:
//   rules:
//     - hosts:
//         - api.example.com
//         - "*.dev.example.com"
//       headers:
//         X-Environment: "development"
//         X-Service: "api-gateway"
//
// Performance considerations:
//   - The plugin processes rules sequentially, so place most specific rules first
//   - Wildcard matching is more expensive than exact matching
//   - Consider rule order for optimal performance in large rule sets
//
// The plugin is compatible with yaegi (Traefik's plugin interpreter) and handles nil next handlers
// gracefully to prevent panics during configuration testing.
package traefik_conditional_headers

import (
	"context"
	"net/http"
	"strings"
)

// Config represents the plugin configuration structure.
// It contains a list of rules that define which headers to set based on hostname matching.
type Config struct {
	// Rules is the ordered list of hostname-to-header mapping rules.
	// Rules are evaluated in sequence, with the first matching rule being applied.
	Rules []Rule `json:"rules"`
}

// Rule defines a mapping between host patterns and HTTP headers to be set.
// Multiple hosts can be specified in a single rule to share the same headers.
type Rule struct {
	// Hosts is the list of hostname patterns that should trigger this rule.
	// Supports exact matches, wildcards (e.g., "*.example.com"), and partial matches.
	Hosts []string `json:"hosts"`

	// Headers is the map of HTTP header names and values to be set when this rule matches.
	// All headers in this map will be set on the matching requests.
	Headers map[string]string `json:"headers"`
}

// CreateConfig creates a new default configuration for the plugin.
// This function is called by Traefik to initialize the plugin with default values.
//
// Returns:
//   - *Config: A pointer to a configuration with empty rules slice.
//
// Example:
//   config := CreateConfig()
//   config.Rules = []Rule{
//       {
//           Hosts: []string{"example.com"},
//           Headers: map[string]string{"X-Custom": "value"},
//       },
//   }
//goland:noinspection GoUnusedExportedFunction
func CreateConfig() *Config {
	return &Config{
		Rules: []Rule{},
	}
}

// conditionalHeaders is the main plugin handler that implements http.Handler interface.
// It contains the plugin state including the next handler in the chain, configured rules,
// and the plugin instance name for debugging purposes.
type conditionalHeaders struct {
	next  http.Handler // The next handler in the Traefik middleware chain
	rules []Rule       // The configured hostname-to-header mapping rules
	name  string       // The plugin instance name for debugging/logging
}

// New creates a new instance of the conditional headers plugin.
// This function is called by Traefik when the plugin is loaded.
//
// Parameters:
//   - ctx: The context for the plugin instance (currently unused but required by interface)
//   - next: The next handler in the Traefik middleware chain
//   - config: The plugin configuration containing rules
//   - name: The name of this plugin instance
//
// Returns:
//   - http.Handler: The plugin handler instance
//   - error: Always returns nil as this plugin cannot fail to initialize
//
// The function gracefully handles nil next handlers to maintain compatibility with yaegi
// (Traefik's plugin interpreter) during configuration testing.
//goland:noinspection GoUnusedExportedFunction,GoUnusedParameter
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &conditionalHeaders{
		next:  next,
		rules: config.Rules,
		name:  name,
	}, nil
}

// ServeHTTP implements the http.Handler interface to process incoming requests.
// This is the main entry point where the plugin evaluates hostname matching rules
// and applies appropriate headers to matching requests.
//
// The algorithm follows these steps:
// 1. Extract the hostname from the incoming request
// 2. Iterate through rules in configuration order
// 3. For each rule, check if any host pattern matches the request hostname
// 4. On first match, apply all headers from that rule and stop processing
// 5. If no rules match, continue to next handler without modifications
//
// Parameters:
//   - responseWriter: The HTTP response writer for writing the response
//   - request: The incoming HTTP request to be processed
//
// Note: This function handles nil next handlers gracefully to maintain compatibility
// with yaegi (Traefik's plugin interpreter) during configuration testing and validation.
func (c *conditionalHeaders) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	host := request.Host

	// Check each rule in order - first match wins
	for _, rule := range c.rules {
		// Check if any of the hosts in this rule match the request hostname
		for _, ruleHost := range rule.Hosts {
			if matchesHost(host, ruleHost) {
				// Apply all headers from this matching rule
				for key, value := range rule.Headers {
					request.Header.Set(key, value)
				}
				// Continue to next handler if available (yaegi compatibility)
				if c.next != nil {
					c.next.ServeHTTP(responseWriter, request)
				}
				return // Stop processing rules after first match
			}
		}
	}

	// No rules matched - continue to next handler without modifications
	// Check if next handler is nil before calling it (needed for yaegi compatibility)
	if c.next != nil {
		c.next.ServeHTTP(responseWriter, request)
	}
}

// matchesHost determines if an incoming request hostname matches a rule host pattern.
// This function supports three matching strategies in order of specificity:
//
// 1. Exact matching: "api.example.com" matches "api.example.com" only
// 2. Wildcard subdomain matching: "*.example.com" matches "api.example.com", "v1.api.example.com", etc.
// 3. Partial string matching: "api" matches "my-api.example.com", "api.example.com", etc.
//
// The function automatically strips port numbers from the incoming host before matching.
//
// Parameters:
//   - incomingHost: The hostname from the incoming HTTP request (may include port)
//   - ruleHost: The hostname pattern from the rule configuration
//
// Returns:
//   - bool: True if the incoming host matches the rule pattern, false otherwise
//
// Matching Examples:
//   - matchesHost("api.example.com:8080", "api.example.com") → true (port stripped)
//   - matchesHost("v1.api.example.com", "*.example.com") → true (wildcard)
//   - matchesHost("my-api.example.com", "api") → true (partial)
//   - matchesHost("other.com", "example.com") → false (no match)
//
// Performance Note: Exact matches are fastest, followed by wildcard matches,
// with partial matches being the most expensive due to string searching.
func matchesHost(incomingHost, ruleHost string) bool {
	// Remove port from incoming host if present (e.g., "example.com:8080" → "example.com")
	if idx := strings.Index(incomingHost, ":"); idx != -1 {
		incomingHost = incomingHost[:idx]
	}

	// 1. Exact match - fastest path
	if incomingHost == ruleHost {
		return true
	}

	// 2. Wildcard subdomain matching (e.g., "*.example.com")
	if strings.HasPrefix(ruleHost, "*.") {
		domain := ruleHost[2:] // Remove "*." prefix to get the base domain
		// Matches if incoming host ends with the domain (with a dot) or equals the domain
		return strings.HasSuffix(incomingHost, "."+domain) || incomingHost == domain
	}

	// 3. Partial string matching (fallback for broader matching)
	// Warning: This can match unexpectedly (e.g., "com" matches "example.com")
	if strings.Contains(incomingHost, ruleHost) {
		return true
	}

	return false
}
