package traefik_conditional_headers

import (
	"context"
	"net/http"
	"strings"
)

type Config struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Hosts   []string          `json:"hosts"` // Multiple hosts
	Headers map[string]string `json:"headers"`
}

//goland:noinspection GoUnusedExportedFunction
func CreateConfig() *Config {
	return &Config{
		Rules: []Rule{},
	}
}

type conditionalHeaders struct {
	next  http.Handler
	rules []Rule
	name  string
}

//goland:noinspection GoUnusedExportedFunction,GoUnusedParameter
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &conditionalHeaders{
		next:  next,
		rules: config.Rules,
		name:  name,
	}, nil
}

func (c *conditionalHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	host := req.Host

	// Check each rule
	for _, rule := range c.rules {
		// Check if any of the hosts in this rule match
		for _, ruleHost := range rule.Hosts {
			if matchesHost(host, ruleHost) {
				// Apply all headers from this rule
				for key, value := range rule.Headers {
					req.Header.Set(key, value)
				}
				c.next.ServeHTTP(rw, req)
				return
			}
		}
	}

	c.next.ServeHTTP(rw, req)
}

// matchesHost checks if the incoming host matches the rule host
// Supports exact match and wildcard subdomain matching
func matchesHost(incomingHost, ruleHost string) bool {
	// Remove port if present
	if idx := strings.Index(incomingHost, ":"); idx != -1 {
		incomingHost = incomingHost[:idx]
	}

	// Exact match
	if incomingHost == ruleHost {
		return true
	}

	// Wildcard match (e.g., *.demo.dev.io)
	if strings.HasPrefix(ruleHost, "*.") {
		domain := ruleHost[2:] // Remove "*."
		return strings.HasSuffix(incomingHost, "."+domain) || incomingHost == domain
	}

	// Contains match (for partial matching)
	if strings.Contains(incomingHost, ruleHost) {
		return true
	}

	return false
}
