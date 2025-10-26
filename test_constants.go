package traefik_conditional_headers

// Test constants to avoid string literal duplication across test files.

// Host constants.
const (
	testHostExampleCom            = "example.com"
	testHostAPIExampleCom         = "api.example.com"
	testHostWildcardExampleCom    = "*.example.com"
	testHostOtherCom              = "other.com"
	testHostTestCom               = "test.com"
	testHostTestExampleCom        = "test.example.com"
	testHostMyAPIExampleCom       = "my-api.example.com"
	testHostV1APIExampleCom       = "v1.api.example.com"
	testHostAPIDevExampleCom      = "api.dev.example.com"
	testHostProductionAPIExample  = "production-api.example.com"
	testHostAuthServiceLocal      = "auth.service.local"
	testHostLoginServiceLocal     = "login.service.local"
	testHostUsersAPIServiceLocal  = "users.api.service.local"
	testHostProfileAPIService     = "profile.api.service.local"
	testHostAdminServiceLocal     = "admin.service.local"
	testHostCacheServiceLocal     = "cache.service.local"
	testHostImagesCDNExample      = "images.cdn.example.com"
	testHostImgCDNExample         = "img.cdn.example.com"
	testHostCDNExample            = "cdn.example.com"
	testHostStaticExample         = "static.example.com"
	testHostAssetsExample         = "assets.example.com"
	testHostExampleComWithPort    = "example.com:8080"
	testHostAPIExampleComWildcard = "*.api.example.com"
	testHostAPISubstring          = "api"
	testHostCOMSubstring          = "com"
	testHostTestSubstring         = "test"
	testHostEmpty                 = ""
	testHostPortOnly              = ":8080"
	testHostWildcardOnly          = "*"
	testHostDotExampleCom         = ".example.com"
	testHostTestAccent            = "tést.example.com"
)

// Host pattern constants.
const (
	testHostWildcardDevExample     = "*.dev.example.com"
	testHostWildcardStagingExample = "*.staging.example.com"
	testHostStagingAPIExample      = "staging-api.example.com"
	testHostWildcardServiceLocal   = "*.service.local"
	testHostWildcardExampleNet     = "*.example.net"
)

// Header name constants.
const (
	testHeaderXTest         = "X-Test"
	testHeaderXCustom       = "X-Custom"
	testHeaderXService      = "X-Service"
	testHeaderXVersion      = "X-Version"
	testHeaderXEnvironment  = "X-Environment"
	testHeaderXDebug        = "X-Debug"
	testHeaderXCORS         = "X-CORS"
	testHeaderXCacheControl = "X-Cache-Control"
	testHeaderXFirst        = "X-First"
	testHeaderXSecond       = "X-Second"
	testHeaderXPortTest     = "X-Port-Test"
	testHeaderXPartial      = "X-Partial"
	testHeaderXServiceType  = "X-Service-Type"
	testHeaderXAuthRequired = "X-Auth-Required"
	testHeaderXPublic       = "X-Public"
	testHeaderXRateLimit    = "X-Rate-Limit"
	testHeaderXRoleRequired = "X-Role-Required"
	testHeaderXAuditLog     = "X-Audit-Log"
	testHeaderXCluster      = "X-Cluster"
	testHeaderXDatacenter   = "X-Datacenter"
	testHeaderXContentType  = "X-Content-Type"
	testHeaderXCompress     = "X-Compress"
	testHeaderXCDN          = "X-CDN"
	testHeaderXLevel        = "X-Level"
	testHeaderXUnicode      = "X-Unicode"
	testHeaderXWildcard     = "X-Wildcard"
	testHeaderXDomain       = "X-Domain"
	testHeaderXRuleID       = "X-Rule-ID"
)

// Header value constants.
const (
	testValueValue             = "value"
	testValueTestValue         = "test-value"
	testValueDevelopment       = "development"
	testValueProduction        = "production"
	testValueTrue              = "true"
	testValueFalse             = "false"
	testValueMatch             = "match"
	testValueSuccess           = "success"
	testValueShouldWin         = "should-win"
	testValueShouldNotWin      = "should-not-win"
	testValueAPI               = "api"
	testValueV1                = "v1"
	testValueV2                = "v2"
	testValueAPIGateway        = "api-gateway"
	testValueV210              = "v2.1.0"
	testValueNoCache           = "no-cache"
	testValueAuthentication    = "authentication"
	testValueUserManagement    = "user-management"
	testValueAdministration    = "administration"
	testValueAdmin             = "admin"
	testValue1000PerHour       = "1000/hour"
	testValue10000PerHour      = "10000/hour"
	testValueProductionCluster = "production"
	testValueUSWest2           = "us-west-2"
	testValueImage             = "image"
	testValueStatic            = "static"
	testValueAPI2              = "api"
	testValueGzip              = "gzip"
	testValueEnabled           = "enabled"
	testValuePublicMaxAge      = "public, max-age=31536000"
	testValuePublicImmutable   = "public, max-age=31536000, immutable"
	testValueSubdomain         = "subdomain"
	testValueAPISubdomain      = "api-subdomain"
	testValueTest              = "test"
	testValuePrimary           = "primary"
)

// Plugin name constants.
const (
	testPluginName            = "test"
	testPluginNameWithDash    = "test-plugin"
	testPluginNameIntegration = "integration-test"
)

// Response constants.
const (
	testResponseOK         = "OK"
	testResponseNextCalled = "Next handler called"
)

// Misc test constants.
const (
	testURLPath = "/test"
)
