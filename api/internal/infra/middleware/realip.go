package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// RealIPConfig holds real IP middleware configuration
type RealIPConfig struct {
	// TrustedProxies is a list of trusted proxy IPs or CIDR ranges
	// If empty, all proxies are trusted
	TrustedProxies []string

	// Headers is a list of headers to check for real IP (in order)
	// Default: ["X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"]
	Headers []string

	// ContextKey is the key used to store the real IP in context
	// Default: "real_ip"
	ContextKey string

	// RecursiveCheck enables recursive checking of X-Forwarded-For
	// to find the first untrusted IP
	RecursiveCheck bool
}

// DefaultRealIPConfig returns default real IP configuration
func DefaultRealIPConfig() RealIPConfig {
	return RealIPConfig{
		TrustedProxies: []string{},
		Headers:        []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"},
		ContextKey:     "real_ip",
		RecursiveCheck: true,
	}
}

// RealIP returns real IP middleware with default config
func RealIP() gin.HandlerFunc {
	return RealIPWithConfig(DefaultRealIPConfig())
}

// RealIPWithConfig returns real IP middleware with custom config
func RealIPWithConfig(cfg RealIPConfig) gin.HandlerFunc {
	if len(cfg.Headers) == 0 {
		cfg.Headers = []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"}
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "real_ip"
	}

	// Parse trusted proxies into CIDR networks
	var trustedNets []*net.IPNet
	for _, proxy := range cfg.TrustedProxies {
		if strings.Contains(proxy, "/") {
			_, network, err := net.ParseCIDR(proxy)
			if err == nil {
				trustedNets = append(trustedNets, network)
			}
		} else {
			// Single IP, convert to /32 or /128
			ip := net.ParseIP(proxy)
			if ip != nil {
				var mask net.IPMask
				if ip.To4() != nil {
					mask = net.CIDRMask(32, 32)
				} else {
					mask = net.CIDRMask(128, 128)
				}
				trustedNets = append(trustedNets, &net.IPNet{IP: ip, Mask: mask})
			}
		}
	}

	return func(c *gin.Context) {
		// Get client IP from remote address
		remoteIP := c.ClientIP()

		// Try to get real IP from headers
		realIP := remoteIP

		for _, header := range cfg.Headers {
			value := c.GetHeader(header)
			if value == "" {
				continue
			}

			// Handle X-Forwarded-For specially (comma-separated list)
			if header == "X-Forwarded-For" {
				ips := strings.Split(value, ",")
				if cfg.RecursiveCheck && len(trustedNets) > 0 {
					// Find first untrusted IP from right to left
					for i := len(ips) - 1; i >= 0; i-- {
						ip := strings.TrimSpace(ips[i])
						if !isTrusted(ip, trustedNets) {
							realIP = ip
							break
						}
					}
				} else {
					// Just use the first (leftmost) IP
					realIP = strings.TrimSpace(ips[0])
				}
				break
			}

			// For other headers, use the value directly
			realIP = strings.TrimSpace(value)
			break
		}

		// Store real IP in context
		c.Set(cfg.ContextKey, realIP)

		c.Next()
	}
}

// isTrusted checks if an IP is in the trusted networks
func isTrusted(ipStr string, trustedNets []*net.IPNet) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, network := range trustedNets {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// GetRealIP retrieves the real IP from context
func GetRealIP(c *gin.Context) string {
	if ip, exists := c.Get("real_ip"); exists {
		return ip.(string)
	}
	return c.ClientIP()
}

// Common headers for different CDN/proxy providers
const (
	HeaderXForwardedFor  = "X-Forwarded-For"
	HeaderXRealIP        = "X-Real-IP"
	HeaderCFConnectingIP = "CF-Connecting-IP" // Cloudflare
	HeaderTrueClientIP   = "True-Client-IP"   // Akamai/Cloudflare Enterprise
	HeaderXClientIP      = "X-Client-IP"      // AWS ELB
)
