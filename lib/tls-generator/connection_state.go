package tlsgenerator

import (
	"crypto/tls"
)

func TlsVersion(state tls.ConnectionState) string {
	switch state.Version {
	case tls.VersionSSL30:
		return "VersionSSL30"
	case tls.VersionTLS10:
		return "VersionTLS10"
	case tls.VersionTLS11:
		return "VersionTLS11"
	case tls.VersionTLS12:
		return "VersionTLS12"
	default:
		return "Unknown TLS version"
	}
}
