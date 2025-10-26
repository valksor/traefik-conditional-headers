# Conditional Headers Plugin for Traefik

A Traefik middleware plugin that sets HTTP request headers conditionally based on the incoming hostname. Supports multiple hosts sharing the same header configuration with flexible matching strategies.

## Features

- **Multiple Hosts per Rule**: Configure several hosts to share the same headers
- **Flexible Matching**: Exact host matching, wildcard subdomain matching, and partial matching
- **Rule Priority**: First matching rule wins for predictable behavior
- **Port Handling**: Automatically strips port numbers from incoming requests
- **Performance Optimized**: Efficient rule processing with configurable order
- **Yaegi Compatible**: Works with Traefik's plugin interpreter for configuration testing

## Installation

This plugin is available from the Traefik Plugin Catalog and can be used without Pilot.

### Option 1: Published Plugin (Recommended)

#### Using Plugin Catalog

**YAML Configuration:**
```yaml
# traefik.yml
experimental:
  plugins:
    conditional-headers:
      moduleName: github.com/valksor/traefik-conditional-headers
      version: "v0.0.3"
```

**CLI Arguments:**
```bash
traefik --experimental.plugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers \
        --experimental.plugins.conditional-headers.version=v0.0.3
```

#### Docker Compose

```yaml
version: '3.8'
services:
  traefik:
    image: traefik:v3.5
    command:
      - "--experimental.plugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers"
      - "--experimental.plugins.conditional-headers.version=v0.0.3"
    volumes:
      - ./traefik.yml:/traefik.yml
    ports:
      - "80:80"
      - "443:443"
```

#### Kubernetes Installation

**Using Helm Values:**
```yaml
# values.yaml
additionalArguments:
  - "--experimental.plugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers"
  - "--experimental.plugins.conditional-headers.version=v0.0.3"
```

**Using ConfigMap:**
```yaml
# traefik-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: traefik-config
data:
  traefik.yml: |
    experimental:
      plugins:
        conditional-headers:
          moduleName: github.com/valksor/traefik-conditional-headers
          version: "v0.0.3"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: traefik
spec:
  template:
    spec:
      containers:
        - name: traefik
          args:
            - "--experimental.plugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers"
            - "--experimental.plugins.conditional-headers.version=v0.0.3"
```

### Development Setup

For developers who want to contribute or modify the plugin:

#### Clone and Test

```bash
# Clone the repository
git clone https://github.com/valksor/traefik-conditional-headers.git
cd traefik-conditional-headers

# Run tests
go test -v ./...

# Run linting
golangci-lint run
```

#### Local Development

Use the source directory directly - no building required:

**YAML Configuration for Local Testing:**
```yaml
# traefik.yml
experimental:
  localPlugins:
    conditional-headers:
      moduleName: github.com/valksor/traefik-conditional-headers
```

**CLI Arguments for Local Testing:**
```bash
traefik --experimental.localPlugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers
```

**Docker with Local Plugin:**
```yaml
version: '3.8'
services:
  traefik:
    image: traefik:v3.5
    command:
      - "--experimental.localPlugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers"
    volumes:
      - ./traefik.yml:/traefik.yml
      - ./traefik-conditional-headers:/plugins-local/src/github.com/valksor/traefik-conditional-headers
    ports:
      - "80:80"
      - "443:443"
```

## Configuration

### Quick Start Example

**YAML Configuration:**
```yaml
http:
  middlewares:
    my-conditional-headers:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - api.example.com
                - "*.dev.example.com"
              headers:
                X-Environment: "development"
                X-Service: "api-gateway"
```

**Docker Compose Labels:**
```yaml
labels:
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].hosts[0]=api.example.com"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].hosts[1]=*.dev.example.com"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].headers.X-Environment=development"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].headers.X-Service=api-gateway"
```

### Multiple Hosts with Same Headers

**YAML Configuration:**
```yaml
http:
  middlewares:
    my-conditional-headers:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - demo.dev.io
                - test.demo.dev.io
                - staging.demo.dev.io
              headers:
                X-App-Kernel-Name: demo.valksor.com
                X-Environment: production
            - hosts:
                - prod.example.com
              headers:
                X-App-Kernel-Name: prod.example.com
                X-Environment: production
```

**CLI Equivalent:**
```bash
traefik \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].hosts[0]=demo.dev.io \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].hosts[1]=test.demo.dev.io \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].hosts[2]=staging.demo.dev.io \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].headers.X-App-Kernel-Name=demo.valksor.com \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].headers.X-Environment=production \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[1].hosts[0]=prod.example.com \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[1].headers.X-App-Kernel-Name=prod.example.com \
  --middlewares.my-conditional-headers.plugin.conditional-headers.rules[1].headers.X-Environment=production
```

### Wildcard Subdomain Matching

**YAML Configuration:**
```yaml
http:
  middlewares:
    my-conditional-headers:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - "*.internal.dev.io"
              headers:
                X-App-Kernel-Name: internal.valksor.com
                X-Network: "internal"
            - hosts:
                - "*.public.dev.io"
              headers:
                X-App-Kernel-Name: public.valksor.com
                X-Network: "public"
```

**Docker Compose Labels:**
```yaml
labels:
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].hosts[0]=*.internal.dev.io"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].headers.X-App-Kernel-Name=internal.valksor.com"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[0].headers.X-Network=internal"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[1].hosts[0]=*.public.dev.io"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[1].headers.X-App-Kernel-Name=public.valksor.com"
  - "traefik.http.middlewares.my-conditional-headers.plugin.conditional-headers.rules[1].headers.X-Network=public"
```

### Real-World Scenarios

#### 1. Microservices Environment Detection

**YAML Configuration:**
```yaml
http:
  middlewares:
    env-detection:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - "*.dev.company.com"
                - "dev-api.company.com"
              headers:
                X-Environment: "development"
                X-Debug: "true"
                X-Log-Level: "debug"
                X-CORS: "*"
            - hosts:
                - "*.staging.company.com"
                - "staging-api.company.com"
              headers:
                X-Environment: "staging"
                X-Debug: "true"
                X-Log-Level: "info"
                X-CORS: "https://company.com"
            - hosts:
                - "*.company.com"
                - "api.company.com"
              headers:
                X-Environment: "production"
                X-Debug: "false"
                X-Log-Level: "warn"
                X-Cache-Control: "no-cache"
```

#### 2. API Versioning and Routing

**YAML Configuration:**
```yaml
http:
  middlewares:
    api-versioning:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - "v1.api.company.com"
                - "api-v1.company.com"
              headers:
                X-API-Version: "v1"
                X-API-Deprecated: "false"
                X-API-Sunset: "2024-12-31"
                X-Rate-Limit: "1000/hour"
            - hosts:
                - "v2.api.company.com"
                - "api.company.com"
              headers:
                X-API-Version: "v2"
                X-API-Latest: "true"
                X-Rate-Limit: "5000/hour"
                X-Features: "advanced,caching,webhooks"
            - hosts:
                - "legacy.api.company.com"
              headers:
                X-API-Version: "legacy"
                X-API-Deprecated: "true"
                X-API-Sunset: "2024-06-30"
                X-Rate-Limit: "100/hour"
```

#### 3. Multi-Tenant SaaS Platform

**YAML Configuration:**
```yaml
http:
  middlewares:
    tenant-routing:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - "*.enterprise.company.com"
                - "company.enterprise.com"
              headers:
                X-Tenant-Type: "enterprise"
                X-Features: "advanced,unlimited,priority-support"
                X-Rate-Limit: "10000/hour"
                X-Storage-Quota: "1TB"
            - hosts:
                - "*.pro.company.com"
                - "company.pro.com"
              headers:
                X-Tenant-Type: "professional"
                X-Features: "standard,api-access,email-support"
                X-Rate-Limit: "5000/hour"
                X-Storage-Quota: "100GB"
            - hosts:
                - "*.company.com"
                - "app.company.com"
              headers:
                X-Tenant-Type: "basic"
                X-Features: "basic,limited"
                X-Rate-Limit: "1000/hour"
                X-Storage-Quota: "10GB"
```

### Advanced Configuration Examples

#### CDN and Static Asset Handling

**YAML Configuration:**
```yaml
http:
  middlewares:
    cdn-headers:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - "images.cdn.company.com"
                - "img.cdn.company.com"
                - "photos.cdn.company.com"
              headers:
                X-Content-Type: "image"
                X-Cache-Control: "public, max-age=31536000, immutable"
                X-Compress: "gzip"
                X-CDN-Provider: "cloudflare"
            - hosts:
                - "static.cdn.company.com"
                - "assets.cdn.company.com"
                - "js.cdn.company.com"
                - "css.cdn.company.com"
              headers:
                X-Content-Type: "static"
                X-Cache-Control: "public, max-age=86400"
                X-Compress: "gzip,br"
                X-CDN-Provider: "fastly"
            - hosts:
                - "api.company.com"
                - "graph.company.com"
              headers:
                X-Content-Type: "api"
                X-Cache-Control: "no-cache, no-store, must-revalidate"
                X-Rate-Limit: "10000/hour"
                X-API-Version: "v2"
```

### Performance Considerations

1. **Rule Order**: Place most specific rules first for better performance
2. **Exact Matches**: Use exact host matches when possible for fastest processing
3. **Wildcards**: Use wildcard patterns sparingly as they're more expensive
4. **Partial Matches**: Avoid partial matches unless necessary, as they're the most expensive

**Optimized Configuration Example:**
```yaml
http:
  middlewares:
    optimized-headers:
      plugin:
        conditional-headers:
          rules:
            # Exact matches first (fastest)
            - hosts:
                - "api.company.com"
                - "admin.company.com"
              headers:
                X-Service: "critical"
            # Wildcard matches second
            - hosts:
                - "*.api.company.com"
                - "*.admin.company.com"
              headers:
                X-Service: "api"
            # Partial matches last (most expensive)
            - hosts:
                - "api"
                - "admin"
              headers:
                X-Service: "fallback"
```

## Troubleshooting

### Common Issues

#### 1. Headers Not Being Applied

**Problem**: Headers are not showing up in your application.

**Solutions**:
- Check that the middleware is properly applied to your router
- Verify the hostname matches exactly (including ports)
- Check rule order - first match wins
- Use Traefik's debug mode to see plugin activity

```bash
# Enable debug logging
traefik --log.level=DEBUG
```

#### 2. Wildcard Rules Not Working

**Problem**: Wildcard patterns like `*.example.com` aren't matching.

**Solutions**:
- Ensure the pattern starts with `*.`
- Check that your base domain is correct
- Remember that `*.example.com` also matches `example.com`

#### 3. Plugin Not Loading

**Problem**: Plugin fails to load or Traefik shows plugin errors.

**Solutions**:
- Verify plugin name is correct: `github.com/valksor/traefik-conditional-headers`
- Check Traefik version compatibility
- Ensure proper plugin configuration format
- Review Traefik logs for specific error messages

### Debug Configuration

**YAML for Debug Mode:**
```yaml
# traefik.yml
log:
  level: DEBUG
accessLog: true

api:
  dashboard: true
  insecure: true
```

**CLI for Debug Mode:**
```bash
traefik --log.level=DEBUG --accesslog=true --api.dashboard=true --api.insecure=true
```

### Testing Headers

Use curl to verify headers are being set correctly:

```bash
# Test specific host
curl -H "Host: api.example.com" http://localhost:8080 -I

# Test wildcard match
curl -H "Host: v1.api.example.com" http://localhost:8080 -I

# Test with custom headers
curl -H "Host: api.example.com" http://localhost:8080 \
  -H "Authorization: Bearer token" -I
```

## Compatibility

### Traefik Version Support

| Plugin Version | Traefik Version | Status | Release Date |
|---------------|-----------------|--------|--------------|
| v0.0.3        | v3.0+           | ✅ Current | 2025-10-26 |
| v0.0.2        | v3.0+           | ⚠️ Legacy | 2025-10-26 |
| v0.0.1        | v3.0+           | ⚠️ Legacy | 2025-10-26 |

### Platform Support

- ✅ Linux (amd64, arm64)
- ✅ macOS (amd64, arm64)
- ✅ Windows (amd64)
- ✅ Docker containers
- ✅ Kubernetes
- ✅ Docker Swarm

### Known Limitations

1. **Header Values**: Only supports string values (no complex objects)
2. **Rule Count**: Performance may degrade with 1000+ rules
3. **Partial Matching**: Can be overly broad - use with caution
4. **Port Handling**: Incoming ports are stripped, but rule patterns should not include ports

## Contributing

We welcome contributions! Here's how to get started:

### Development Setup

1. **Prerequisites**:
   - Go 1.24+
   - Docker
   - Traefik v3.0+

2. **Clone and Build**:
   ```bash
   git clone https://github.com/valksor/traefik-conditional-headers.git
   cd traefik-conditional-headers
   go mod tidy
   ```

3. **Run Tests**:
   ```bash
   # Run unit tests
   go test -v ./...

   # Run integration tests
   go test -v -tags=integration ./...

   # Run benchmarks
   go test -bench=. -benchmem ./...

   # Check code quality
   go vet ./...
   go fmt ./...
   ```

4. **Local Testing**:
   ```bash
   # Start test Traefik instance
   docker run -it --rm \
     -v $(pwd):/plugins-local \
     -p 8080:8080 \
     traefik:v3.5 \
     --experimental.localPlugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers
   ```

### Submitting Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** and ensure:
   - All tests pass
   - Code follows Go conventions
   - Documentation is updated
   - Examples are tested

3. **Submit a pull request** with:
   - Clear description of changes
   - Test cases for new functionality
   - Updated documentation if needed
   - Examples of usage

### Code Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Add comprehensive tests for new features
- Update documentation for any API changes
- Use meaningful commit messages
- Keep backwards compatibility when possible

### Testing Guidelines

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test complete request flows
- **Edge Cases**: Test unusual hostnames and configurations
- **Performance**: Benchmark with realistic rule sets
- **Compatibility**: Test with different Traefik versions

### Release Process

This project follows semantic versioning and has released 3 versions:

1. **v0.0.1** (2025-10-26): Initial release with core functionality
2. **v0.0.2** (2025-10-26): Configuration file naming update
3. **v0.0.3** (2025-10-26): Enhanced CI/CD and stability improvements

Future releases will follow this process:
1. Finalize features and testing
2. Update version in documentation
3. Update CHANGELOG.md with release notes
4. Create git tag: `git tag v0.0.x`
5. Push to repository
6. Create GitHub release
7. Publish to Traefik plugin catalog
8. Update documentation with new version information

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://github.com/valksor/traefik-conditional-headers)
- 🐛 [Issue Tracker](https://github.com/valksor/traefik-conditional-headers/issues)
- 💬 [Discussions](https://github.com/orgs/valksor/discussions)

## Related Plugins

- [Traefik Plugin Catalog](https://plugins.traefik.io/)
- [Custom Headers Plugin](https://github.com/traefik/traefik/tree/master/plugins/pilot/examples)
- [Rate Limit Plugin](https://github.com/traefik/traefik-ratelimit)
