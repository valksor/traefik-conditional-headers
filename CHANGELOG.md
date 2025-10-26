# Changelog

All notable changes to the Traefik Conditional Headers Plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Enhanced Documentation**: Comprehensive documentation with real-world examples, installation guides, and troubleshooting
- **Test Coverage**: Enhanced test data with comprehensive scenarios and edge cases
- **Code Documentation**: Complete godoc documentation for all exported functions and types

## [0.0.3] - 2025-10-26

### Added
- **CI/CD Pipeline**: Cross-platform Go builds for wider compatibility
- **Code Quality**: Automated linting and testing with golangci-lint
- **Logo**: Added plugin logo to middleware manifest
- **Makefile**: Build automation for development and release processes

### Improved
- **Stability**: Enhanced request handling and error handling
- **Yaegi Compatibility**: Addressed potential nil pointer dereference for yaegi compatibility
- **Test Suite**: Improved integration tests and edge case coverage
- **Code Style**: Consistent code formatting and structure

### Fixed
- **Error Handling**: Better handling of edge cases and nil pointers
- **Test Reliability**: More robust and comprehensive test suite

## [0.0.2] - 2025-10-26

### Changed
- **Configuration File**: Renamed `.traefik.yaml` to `.traefik.yml` for consistency with Traefik conventions

## [0.0.1] - 2025-10-26

### Added
- **Initial Release**: First release of the Conditional Headers Plugin
- **Core Functionality**:
  - Multiple host support per rule
  - Exact hostname matching
  - Wildcard subdomain matching (e.g., `*.example.com`)
  - Partial string matching for flexible patterns
- **First Match Algorithm**: Rules are evaluated in order with first match winning
- **Port Handling**: Automatic port stripping from incoming hostnames
- **Comprehensive Testing**:
  - Unit tests for all core functions
  - Integration tests for complete request flows
  - Edge case testing and performance benchmarks
- **Traefik Integration**: Full compatibility with Traefik v2.8+ middleware system
- **Documentation**: Basic usage examples and configuration reference

## Version Support Matrix

| Version | Traefik Support | Status | Release Date |
|---------|-----------------|--------|--------------|
| 0.0.3   | v3.0+           | ✅ Current | 2025-10-26 |
| 0.0.2   | v3.0+           | ⚠️ Legacy | 2025-10-26 |
| 0.0.1   | v3.0+           | ⚠️ Legacy | 2025-10-26 |

### Platform Support

All versions support:
- ✅ Linux (amd64, arm64)
- ✅ macOS (amd64, arm64)
- ✅ Windows (amd64)
- ✅ Docker containers
- ✅ Kubernetes
- ✅ Docker Swarm

### Upgrade Guide

#### From v0.0.1/v0.0.2 to v0.0.3
No breaking changes. Simply update your configuration:

```yaml
# Update the version in your configuration
pilot:
  plugins:
    conditional-headers:
      moduleName: github.com/valksor/traefik-conditional-headers
      version: "v0.0.3"
```

### Known Issues
- No known issues in current version

### Contributing
Contributions are welcome! Please see the [Contributing Guidelines](README.md#contributing) for details.

### Support
- 📖 [Documentation](README.md)
- 🐛 [Issue Tracker](https://github.com/valksor/traefik-conditional-headers/issues)
- 💬 [Discussions](https://github.com/orgs/valksor/discussions)
