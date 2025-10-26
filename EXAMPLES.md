# Real-World Configuration Examples

This document provides detailed, production-ready configuration examples for the Traefik Conditional Headers plugin in various scenarios.

## Table of Contents

- [E-commerce Platform](#e-commerce-platform)
- [SaaS Application](#saas-application)
- [Microservices Architecture](#microservices-architecture)
- [Content Delivery Network](#content-delivery-network)
- [API Gateway](#api-gateway)
- [Multi-Region Deployment](#multi-region-deployment)
- [Educational Platform](#educational-platform)
- [Healthcare System](#healthcare-system)

---

## E-commerce Platform

### Scenario
An online retail store with separate services for product catalog, user accounts, payment processing, and admin interface.

### Architecture
- **Product API**: `catalog.shop.com`, `api.products.shop.com`
- **User Service**: `accounts.shop.com`, `profile.shop.com`
- **Payment Gateway**: `payment.shop.com`, `checkout.shop.com`
- **Admin Panel**: `admin.shop.com`, `staff.shop.com`
- **Static Assets**: `cdn.shop.com`, `assets.shop.com`

### Configuration

```yaml
# traefik.yml
http:
  middlewares:
    shop-headers:
      plugin:
        conditional-headers:
          rules:
            # Product catalog services
            - hosts:
                - "catalog.shop.com"
                - "api.products.shop.com"
                - "*.catalog.shop.com"
              headers:
                X-Service: "product-catalog"
                X-API-Version: "v2"
                X-Cache-Control: "public, max-age=3600"
                X-Rate-Limit: "10000/hour"
                X-Features: "search,recommendations,inventory"

            # User account services
            - hosts:
                - "accounts.shop.com"
                - "profile.shop.com"
                - "users.shop.com"
              headers:
                X-Service: "user-management"
                X-API-Version: "v1"
                X-Auth-Required: "true"
                X-GDPR: "compliant"
                X-Data-Encryption: "AES-256"
                X-Rate-Limit: "5000/hour"

            # Payment processing (highest security)
            - hosts:
                - "payment.shop.com"
                - "checkout.shop.com"
                - "billing.shop.com"
              headers:
                X-Service: "payment-gateway"
                X-Security: "maximum"
                X-PCI-DSS: "compliant"
                X-Encryption: "TLS-1.3"
                X-Audit: "enabled"
                X-Auth-Required: "true"
                X-MFA: "required"
                X-Rate-Limit: "1000/hour"
                X-Timeout: "30"

            # Administrative interfaces
            - hosts:
                - "admin.shop.com"
                - "staff.shop.com"
                - "internal.shop.com"
              headers:
                X-Service: "administration"
                X-Auth-Required: "true"
                X-Role-Required: "admin|staff"
                X-Audit-Log: "true"
                X-IP-Whitelist: "enabled"
                X-Session-Timeout: "1800"
                X-Security: "high"

            # Static assets and CDN
            - hosts:
                - "cdn.shop.com"
                - "assets.shop.com"
                - "images.shop.com"
                - "*.cdn.shop.com"
              headers:
                X-Service: "static-content"
                X-Cache-Control: "public, max-age=31536000, immutable"
                X-Compress: "gzip,br"
                X-CDN: "cloudflare"
                X-Brotli: "enabled"
                X-Storage: "s3-intelligent-tiering"

            # Mobile API endpoints
            - hosts:
                - "mobile-api.shop.com"
                - "app.shop.com"
                - "m.shop.com"
              headers:
                X-Service: "mobile-api"
                X-Platform: "mobile"
                X-API-Version: "v3"
                X-Push-Notifications: "enabled"
                X-Offline-Support: "enabled"
                X-Rate-Limit: "20000/hour"

            # Development environment
            - hosts:
                - "*.dev.shop.com"
                - "dev-shop.company.internal"
              headers:
                X-Environment: "development"
                X-Debug: "true"
                X-Log-Level: "debug"
                X-CORS: "*"
                X-Source-Maps: "enabled"
                X-Hot-Reload: "true"
                X-Database: "dev-postgres"
```

### Docker Compose Implementation

```yaml
# docker-compose.yml
version: '3.8'
services:
  traefik:
    image: traefik:v3.5
    command:
      - "--experimental.plugins.conditional-headers.moduleName=github.com/valksor/traefik-conditional-headers"
      - "--experimental.plugins.conditional-headers.version=v1.2.0"
    volumes:
      - ./traefik.yml:/traefik.yml
    ports:
      - "80:80"
      - "443:443"
    labels:
      # E-commerce platform headers
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[0].hosts[0]=catalog.shop.com"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[0].hosts[1]=api.products.shop.com"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[0].headers.X-Service=product-catalog"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[0].headers.X-API-Version=v2"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[0].headers.X-Rate-Limit=10000/hour"

      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[1].hosts[0]=payment.shop.com"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[1].headers.X-Security=maximum"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[1].headers.X-PCI-DSS=compliant"
      - "traefik.http.middlewares.shop-headers.plugin.conditional-headers.rules[1].headers.X-MFA=required"
```

---

## SaaS Application

### Scenario
A multi-tenant SaaS platform with different subscription tiers and feature sets.

### Architecture
- **Enterprise Tier**: Custom subdomains, full feature access
- **Professional Tier**: Standard features, API access
- **Basic Tier**: Limited features, community support
- **Public Marketing**: Marketing website, documentation

### Configuration

```yaml
http:
  middlewares:
    saas-headers:
      plugin:
        conditional-headers:
          rules:
            # Enterprise customers
            - hosts:
                - "*.enterprise.saasapp.com"
                - "company.enterprise.saasapp.com"
              headers:
                X-Tenant-Type: "enterprise"
                X-Features: "advanced,unlimited,priority-support,custom-integrations,white-label,advanced-analytics"
                X-Rate-Limit: "100000/hour"
                X-Storage-Quota: "unlimited"
                X-Support-Level: "24/7-dedicated"
                X-SLA: "99.99%"
                X-API-Access: "full"
                X-Custom-Domain: "enabled"
                X-SSO: "enabled"
                X-Audit-Retention: "10-years"

            # Professional tier customers
            - hosts:
                - "*.pro.saasapp.com"
                - "company.pro.saasapp.com"
                - "app.pro.saasapp.com"
              headers:
                X-Tenant-Type: "professional"
                X-Features: "standard,api-access,email-support,advanced-analytics,team-collaboration,automated-backups"
                X-Rate-Limit: "50000/hour"
                X-Storage-Quota: "1TB"
                X-Support-Level: "business-hours"
                X-SLA: "99.9%"
                X-API-Access: "standard"
                X-Team-Members: "50"
                X-Integrations: "slack,teams,github,jira"

            # Basic tier customers
            - hosts:
                - "*.saasapp.com"
                - "app.saasapp.com"
                - "my.saasapp.com"
              headers:
                X-Tenant-Type: "basic"
                X-Features: "basic,limited,community-support,email-support"
                X-Rate-Limit: "10000/hour"
                X-Storage-Quota: "100GB"
                X-Support-Level: "community"
                X-SLA: "99%"
                X-API-Access: "limited"
                X-Team-Members: "5"
                X-Integrations: "email,calendar"

            # API gateway
            - hosts:
                - "api.saasapp.com"
                - "graph.saasapp.com"
                - "webhook.saasapp.com"
              headers:
                X-Service: "api-gateway"
                X-API-Version: "v2"
                X-Rate-Limit: "1000000/hour"
                X-Authentication: "oauth2,jwt,apikey"
                X-Monitoring: "prometheus,datadog"
                X-Documentation: "https://docs.saasapp.com"

            # Developer portal
            - hosts:
                - "developers.saasapp.com"
                - "docs.saasapp.com"
                - "api-docs.saasapp.com"
              headers:
                X-Service: "developer-portal"
                X-Content-Type: "documentation"
                X-Cache-Control: "public, max-age=3600"
                X-Analytics: "google-analytics"
                X-Search: "algolia"

            # Admin and internal tools
            - hosts:
                - "admin.saasapp.com"
                - "internal.saasapp.com"
                - "ops.saasapp.com"
              headers:
                X-Service: "internal-tools"
                X-Auth-Required: "true"
                X-Role-Required: "admin|ops"
                X-Audit-Log: "true"
                X-Security: "maximum"
                X-VPN-Required: "true"
                X-Session-Timeout: "3600"

            # Status and monitoring
            - hosts:
                - "status.saasapp.com"
                - "health.saasapp.com"
                - "uptime.saasapp.com"
              headers:
                X-Service: "status-monitoring"
                X-Public: "true"
                X-Refresh-Interval: "60"
                X-Incident-Communication: "statuspage"
```

---

## Microservices Architecture

### Scenario
A microservices-based application with dozens of services communicating through an API gateway.

### Configuration

```yaml
http:
  middlewares:
    microservice-headers:
      plugin:
        conditional-headers:
          rules:
            # Core services
            - hosts:
                - "auth.service.local"
                - "identity.service.local"
                - "sso.service.local"
              headers:
                X-Service-Type: "authentication"
                X-Criticality: "critical"
                X-Dependencies: "database,redis,email-service"
                X-Monitoring: "prometheus,grafana"
                X-Tracing: "jaeger"
                X-Circuit-Breaker: "enabled"
                X-Health-Check: "/health"

            # User services
            - hosts:
                - "users.service.local"
                - "profile.service.local"
                - "notification.service.local"
              headers:
                X-Service-Type: "user-management"
                X-Criticality: "high"
                X-Dependencies: "auth-service,database,notification-queue"
                X-Cache: "redis"
                X-Data-Privacy: "gdpr-compliant"

            # Business logic services
            - hosts:
                - "order.service.local"
                - "payment.service.local"
                - "inventory.service.local"
                - "shipping.service.local"
              headers:
                X-Service-Type: "business-logic"
                X-Criticality: "critical"
                X-Dependencies: "database,queue,payment-gateway"
                X-Transactional: "true"
                X-Audit: "enabled"

            # Analytics and reporting
            - hosts:
                - "analytics.service.local"
                - "reporting.service.local"
                - "data-warehouse.service.local"
              headers:
                X-Service-Type: "analytics"
                X-Criticality: "medium"
                X-Dependencies: "data-lake,redis"
                X-Batch-Processing: "true"
                X-Retention: "7-years"

            # External API integrations
            - hosts:
                - "integrations.service.local"
                - "webhooks.service.local"
                - "external-apis.service.local"
              headers:
                X-Service-Type: "integration"
                X-Criticality: "high"
                X-Dependencies: "redis,queue"
                X-Rate-Limiting: "external"
                X-Timeout: "30"

            # Development environment
            - hosts:
                - "*.dev.service.local"
                - "dev-gateway.service.local"
              headers:
                X-Environment: "development"
                X-Debug: "true"
                X-Log-Level: "debug"
                X-Profiling: "enabled"
                X-Hot-Reload: "true"

            # Testing environment
            - hosts:
                - "*.test.service.local"
                - "test-gateway.service.local"
              headers:
                X-Environment: "testing"
                X-Test-Mode: "true"
                X-Mock-External: "true"
                X-Database: "test-db"
                X-Queue: "test-queue"
```

---

## Content Delivery Network

### Scenario
A CDN setup with different caching strategies for various types of content.

### Configuration

```yaml
http:
  middlewares:
    cdn-headers:
      plugin:
        conditional-headers:
          rules:
            # Images and media files
            - hosts:
                - "images.cdn.example.com"
                - "media.cdn.example.com"
                - "photos.cdn.example.com"
                - "videos.cdn.example.com"
              headers:
                X-Content-Type: "media"
                X-Cache-Control: "public, max-age=31536000, immutable"
                X-Compress: "gzip,br"
                X-CDN-Provider: "cloudflare"
                X-Storage: "s3-intelligent-tiering"
                X-Optimization: "auto-compression,format-conversion"
                X-Watermark: "enabled"

            # Static assets (CSS, JS)
            - hosts:
                - "static.cdn.example.com"
                - "assets.cdn.example.com"
                - "js.cdn.example.com"
                - "css.cdn.example.com"
                - "fonts.cdn.example.com"
              headers:
                X-Content-Type: "static"
                X-Cache-Control: "public, max-age=86400"
                X-Compress: "gzip,br"
                X-CDN-Provider: "fastly"
                X-Brotli: "enabled"
                X-Minification: "enabled"
                X-Bundling: "enabled"
                X-Subresource-Integrity: "enabled"

            # API responses
            - hosts:
                - "api.cdn.example.com"
                - "graph.cdn.example.com"
              headers:
                X-Content-Type: "api"
                X-Cache-Control: "public, max-age=300"
                X-CDN-Provider: "cloudflare"
                X-Edge-Computing: "enabled"
                X-Rate-Limit: "10000/hour"
                X-Authentication: "jwt,apikey"

            # Dynamic content
            - hosts:
                - "dynamic.cdn.example.com"
                - "personalized.cdn.example.com"
              headers:
                X-Content-Type: "dynamic"
                X-Cache-Control: "private, max-age=60"
                X-Personalization: "enabled"
                X-A-B-Testing: "enabled"
                X-Cookie-Handling: "respect"
                X-Geo-Targeting: "enabled"

            # File downloads
            - hosts:
                - "download.cdn.example.com"
                - "files.cdn.example.com"
                - "releases.cdn.example.com"
              headers:
                X-Content-Type: "download"
                X-Cache-Control: "public, max-age=3600"
                X-Content-Disposition: "attachment"
                X-Virus-Scan: "enabled"
                X-Bandwidth-Limit: "10MB/s"
                X-Download-Tracking: "enabled"

            # Edge computing functions
            - hosts:
                - "edge.cdn.example.com"
                - "compute.cdn.example.com"
              headers:
                X-Content-Type: "edge-computing"
                X-Cache-Control: "no-cache"
                X-Edge-Functions: "enabled"
                X-Serverless: "true"
                X-Compute-Limit: "100ms"
                X-Memory-Limit: "128MB"
```

---

## API Gateway

### Scenario
An API gateway managing multiple backend services with different authentication and rate limiting requirements.

### Configuration

```yaml
http:
  middlewares:
    api-gateway-headers:
      plugin:
        conditional-headers:
          rules:
            # Public API v1 (deprecated)
            - hosts:
                - "v1.api.example.com"
                - "api-v1.example.com"
                - "legacy.api.example.com"
              headers:
                X-API-Version: "v1"
                X-API-Deprecated: "true"
                X-API-Sunset: "2024-06-30"
                X-Rate-Limit: "1000/hour"
                X-Authentication: "apikey"
                X-Features: "basic,read-only"
                X-Support: "community"
                X-Documentation: "https://docs.example.com/api/v1"

            # Public API v2 (current stable)
            - hosts:
                - "api.example.com"
                - "v2.api.example.com"
                - "graph.example.com"
              headers:
                X-API-Version: "v2"
                X-API-Stable: "true"
                X-Rate-Limit: "10000/hour"
                X-Authentication: "oauth2,jwt,apikey"
                X-Features: "advanced,caching,webhooks,real-time"
                X-Support: "professional"
                X-Documentation: "https://docs.example.com/api/v2"

            # API v3 (beta)
            - hosts:
                - "v3.api.example.com"
                - "beta.api.example.com"
                - "next.api.example.com"
              headers:
                X-API-Version: "v3"
                X-API-Stability: "beta"
                X-Rate-Limit: "50000/hour"
                X-Authentication: "oauth2,jwt"
                X-Features: "experimental,ai,machine-learning"
                X-Support: "early-adopters"
                X-Feedback: "encouraged"
                X-Documentation: "https://docs.example.com/api/v3"

            # Partner API
            - hosts:
                - "partners.api.example.com"
                - "integration.api.example.com"
              headers:
                X-API-Version: "v2"
                X-API-Type: "partner"
                X-Rate-Limit: "100000/hour"
                X-Authentication: "mutual-tls,oauth2"
                X-SLA: "99.9%"
                X-Support: "dedicated"
                X-Contract: "required"
                X-Webhooks: "enabled"

            # Internal API
            - hosts:
                - "internal.api.example.com"
                - "service-api.example.com"
                - "microservice-api.example.com"
              headers:
                X-API-Version: "v2"
                X-API-Type: "internal"
                X-Rate-Limit: "unlimited"
                X-Authentication: "service-account,mTLS"
                X-Network: "internal-vpc"
                X-Monitoring: "comprehensive"
                X-Tracing: "enabled"
                X-Debug: "true"

            # Admin API
            - hosts:
                - "admin.api.example.com"
                - "management.api.example.com"
              headers:
                X-API-Version: "v2"
                X-API-Type: "administrative"
                X-Rate-Limit: "5000/hour"
                X-Authentication: "mfa,oauth2"
                X-Audit: "enabled"
                X-Role-Required: "admin|operator"
                X-Security: "maximum"
                X-Session-Timeout: "1800"

            # WebSocket endpoints
            - hosts:
                - "ws.api.example.com"
                - "websocket.api.example.com"
                - "realtime.api.example.com"
              headers:
                X-API-Type: "websocket"
                X-Connection: "upgrade"
                X-Protocol: "websocket"
                X-Rate-Limit: "1000/minute"
                X-Authentication: "jwt"
                X-Message-Size: "1MB"
                X-Compression: "enabled"
```

---

## Multi-Region Deployment

### Scenario
A globally distributed application with region-specific routing and compliance requirements.

### Configuration

```yaml
http:
  middlewares:
    regional-headers:
      plugin:
        conditional-headers:
          rules:
            # US East region
            - hosts:
                - "us-east.api.example.com"
                - "use1.api.example.com"
                - "virginia.api.example.com"
              headers:
                X-Region: "us-east-1"
                X-Availability-Zone: "use1-az1,use1-az2,use1-az3"
                X-Compliance: "SOC2,HIPAA,GDPR"
                X-Data-Residency: "US"
                X-Latency: "<50ms"
                X-Backup: "cross-region"
                X-Monitoring: "us-central"

            # US West region
            - hosts:
                - "us-west.api.example.com"
                - "usw2.api.example.com"
                - "california.api.example.com"
              headers:
                X-Region: "us-west-2"
                X-Availability-Zone: "usw2-az1,usw2-az2"
                X-Compliance: "CCPA,CPRA"
                X-Data-Residency: "US"
                X-Latency: "<30ms"
                X-Backup: "cross-region"
                X-Optimization: "us-west-population"

            # Europe region
            - hosts:
                - "eu-west.api.example.com"
                - "euw1.api.example.com"
                - "ireland.api.example.com"
              headers:
                X-Region: "eu-west-1"
                X-Availability-Zone: "euw1-az1,euw1-az2"
                X-Compliance: "GDPR,ePrivacy,Schrems-II"
                X-Data-Residency: "EU"
                X-Latency: "<40ms"
                X-Backup: "cross-region"
                X-Language: "en,fr,de,es,it"

            # Asia Pacific region
            - hosts:
                - "ap-southeast.api.example.com"
                - "aps1.api.example.com"
                - "singapore.api.example.com"
              headers:
                X-Region: "ap-southeast-1"
                X-Availability-Zone: "aps1-az1,aps1-az2"
                X-Compliance: "PDPA,PDPO"
                X-Data-Residency: "Singapore"
                X-Latency: "<60ms"
                X-Backup: "cross-region"
                X-Language: "en,zh,ja,ko,th,vn"

            # Global CDN
            - hosts:
                - "global.cdn.example.com"
                - "*.cdn.example.com"
              headers:
                X-Scope: "global"
                X-Cache-Strategy: "aggressive"
                X-Edge-Locations: "200+"
                X-Optimization: "geo-routing,compression,optimization"
                X-Fallback: "regional"
                X-Monitoring: "global-dashboard"

            # Regional compliance endpoints
            - hosts:
                - "compliance.api.example.com"
                - "gdpr.api.example.com"
                - "privacy.api.example.com"
              headers:
                X-Service: "compliance"
                X-Audit-Trail: "enabled"
                X-Data-Retention: "7-years"
                X-Right-to-Delete: "enabled"
                X-Data-Portability: "enabled"
                X-Consent-Management: "enabled"
                X-Encryption: "AES-256,TLS-1.3"
```

---

## Educational Platform

### Scenario
An online learning platform with different user roles and content types.

### Configuration

```yaml
http:
  middlewares:
    education-headers:
      plugin:
        conditional-headers:
          rules:
            # Student portal
            - hosts:
                - "students.education.com"
                - "learn.education.com"
                - "courses.education.com"
              headers:
                X-User-Role: "student"
                X-Features: "course-access,progress-tracking,submissions"
                X-Rate-Limit: "1000/hour"
                X-Content-Type: "educational"
                X-Learning-Analytics: "enabled"
                X-Accessibility: "WCAG-2.1"

            # Instructor portal
            - hosts:
                - "instructors.education.com"
                - "teach.education.com"
                - "faculty.education.com"
              headers:
                X-User-Role: "instructor"
                X-Features: "course-creation,grading,analytics,communication"
                X-Rate-Limit: "5000/hour"
                X-Content-Type: "educational-admin"
                X-Grade-Book: "enabled"
                X-Student-Management: "enabled"

            # Administration
            - hosts:
                - "admin.education.com"
                - "administration.education.com"
                - "system.education.com"
              headers:
                X-User-Role: "administrator"
                X-Features: "full-access,user-management,system-configuration"
                X-Rate-Limit: "10000/hour"
                X-Audit: "enabled"
                X-Security: "high"
                X-Data-Privacy: "FERPA-compliant"

            # Video content delivery
            - hosts:
                - "video.education.com"
                - "lectures.education.com"
                - "media.education.com"
              headers:
                X-Content-Type: "video-educational"
                X-Cache-Control: "public, max-age=86400"
                X-Streaming: "adaptive-bitrate"
                X-Captions: "auto-generated"
                X-Transcripts: "enabled"
                X-Offline-Download: "premium"

            # Interactive content
            - hosts:
                - "interactive.education.com"
                - "labs.education.com"
                - "simulations.education.com"
              headers:
                X-Content-Type: "interactive-educational"
                X-Rate-Limit: "500/hour"
                X-Real-Time: "enabled"
                X-Collaboration: "enabled"
                X-Save-Progress: "enabled"
                X-Performance: "optimized"

            # Assessment and testing
            - hosts:
                - "exams.education.com"
                - "testing.education.com"
                - "assessment.education.com"
              headers:
                X-Service: "assessment"
                X-Security: "maximum"
                X-Timer: "strict"
                X-Proctoring: "enabled"
                X-Audit: "comprehensive"
                X-Session-Lock: "enabled"
                X-Cheating-Detection: "ai-powered"

            # Library and resources
            - hosts:
                - "library.education.com"
                - "resources.education.com"
                - "research.education.com"
              headers:
                X-Service: "library"
                X-Content-Type: "educational-resources"
                X-Cache-Control: "public, max-age=3600"
                X-Search: "full-text"
                X-Citations: "auto-format"
                X-Plagiarism-Check: "enabled"
```

---

## Healthcare System

### Scenario
A healthcare platform with strict compliance and security requirements.

### Configuration

```yaml
http:
  middlewares:
    healthcare-headers:
      plugin:
        conditional-headers:
          rules:
            # Patient portal (HIPAA compliant)
            - hosts:
                - "patients.healthcare.com"
                - "mychart.healthcare.com"
                - "portal.healthcare.com"
              headers:
                X-Service: "patient-portal"
                X-Compliance: "HIPAA,HITECH"
                X-Data-Encryption: "AES-256"
                X-Authentication: "mfa-required"
                X-Audit: "comprehensive"
                X-Session-Timeout: "900"
                X-Data-Residency: "US"
                X-Access-Control: "role-based"
                X-Privacy: "maximum"

            # Provider portal
            - hosts:
                - "providers.healthcare.com"
                - "doctors.healthcare.com"
                - "clinical.healthcare.com"
              headers:
                X-Service: "provider-portal"
                X-User-Role: "healthcare-provider"
                X-Compliance: "HIPAA,HITECH,21CFR11"
                X-Authentication: "mfa-required"
                X-Audit: "comprehensive"
                X-EMR-Integration: "epic,cerner"
                X-Clinical-Decision: "support"
                X-Prescription: "e-prescribe"

            # Emergency services
            - hosts:
                - "emergency.healthcare.com"
                - "er.healthcare.com"
                - "urgent.healthcare.com"
              headers:
                X-Service: "emergency-services"
                X-Priority: "critical"
                X-Availability: "99.999%"
                X-Response-Time: "<1s"
                X-Authentication: "fast-access"
                X-Audit: "emergency-priority"
                X-Notification: "real-time"
                X-Redundancy: "multi-region"

            # Medical imaging
            - hosts:
                - "imaging.healthcare.com"
                - "radiology.healthcare.com"
                - "dicom.healthcare.com"
              headers:
                X-Service: "medical-imaging"
                X-Compliance: "HIPAA,DICOM"
                X-Data-Type: "medical-images"
                X-Storage: "pacs-integrated"
                X-Compression: "lossless"
                X-Viewing: "diagnostic-quality"
                X-Sharing: "secure-link"
                X-Backup: "geo-redundant"

            # Laboratory services
            - hosts:
                - "lab.healthcare.com"
                - "laboratory.healthcare.com"
                - "pathology.healthcare.com"
              headers:
                X-Service: "laboratory"
                X-Compliance: "HIPAA,CLIA,CAP"
                X-Data-Type: "lab-results"
                X-Accuracy: "validated"
                X-Turnaround: "24-48h"
                X-QC: "continuous"
                X-Integration: "lis,emr"
                X-Reporting: "structured"

            # Pharmacy and medication
            - hosts:
                - "pharmacy.healthcare.com"
                - "medication.healthcare.com"
                - "prescription.healthcare.com"
              headers:
                X-Service: "pharmacy"
                X-Compliance: "HIPAA,DEA,21CFR11"
                X-Authentication: "multi-factor"
                X-Prescription: "e-prescribe"
                X-Drug-Interaction: "checking"
                X-Allergy-Alert: "enabled"
                X-Inventory: "real-time"
                X-Counseling: "automated"

            # Research and clinical trials
            - hosts:
                - "research.healthcare.com"
                - "trials.healthcare.com"
                - "clinical-research.healthcare.com"
              headers:
                X-Service: "research"
                X-Compliance: "HIPAA,GCP,21CFR11"
                X-Data-Anonymization: "required"
                X-IRB: "approved"
                X-Consent: "digital"
                X-Data-Retention: "15-years"
                X-Sharing: "de-identified"
                X-Analytics: "research-only"

            # Telehealth services
            - hosts:
                - "telehealth.healthcare.com"
                - "virtual.healthcare.com"
                - "video.healthcare.com"
              headers:
                X-Service: "telehealth"
                X-Compliance: "HIPAA,HITECH"
                X-Encryption: "end-to-end"
                X-Video-Quality: "HD"
                X-Recording: "patient-consent"
                X-Platform: "HIPAA-compliant"
                X-Billing: "telehealth-codes"
                X-Documentation: "auto-generated"
```

---

## Testing and Validation

### Testing Configurations

Use these curl commands to validate your configurations:

```bash
# Test e-commerce platform
curl -H "Host: catalog.shop.com" http://localhost:8080 -I
curl -H "Host: payment.shop.com" http://localhost:8080 -I

# Test SaaS application
curl -H "Host: company.enterprise.saasapp.com" http://localhost:8080 -I
curl -H "Host: app.saasapp.com" http://localhost:8080 -I

# Test API gateway
curl -H "Host: api.example.com" http://localhost:8080 -I
curl -H "Host: v3.api.example.com" http://localhost:8080 -I

# Test healthcare system
curl -H "Host: patients.healthcare.com" http://localhost:8080 -I
curl -H "Host: emergency.healthcare.com" http://localhost:8080 -I
```

### Performance Monitoring

Monitor header performance with:

```bash
# Benchmark with multiple requests
for i in {1..1000}; do
  curl -H "Host: api.example.com" http://localhost:8080 -o /dev/null -s -w "%{time_total}\n"
done
```

---

## Security Considerations

When implementing these configurations in production:

1. **Authentication**: Always validate headers don't expose sensitive information
2. **Rate Limiting**: Implement appropriate rate limiting for each service tier
3. **HTTPS**: Ensure all endpoints use HTTPS in production
4. **Audit Logging**: Enable comprehensive audit trails for sensitive services
5. **Compliance**: Verify compliance with relevant regulations (GDPR, HIPAA, PCI-DSS)
6. **Regular Updates**: Keep configurations updated with security best practices

---

## Need Help?

For questions about these examples or to request additional scenarios:

- 📖 [Main Documentation](README.md)
- 🐛 [Issue Tracker](https://github.com/valksor/traefik-conditional-headers/issues)
- 💬 [Discussions](https://github.com/valksor/traefik-conditional-headers/discussions)
- 📧 [Email Support](mailto:support@valksor.com)