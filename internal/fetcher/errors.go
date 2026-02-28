package fetcher

import (
	"errors"
	"net"
	"strings"
)

func ClassifyError(err error) string {
	if err == nil {
		return ""
	}

	// context timeout or deadline exceeded
	if errors.Is(err, net.ErrClosed) {
		return "connection_closed"
	}

	// check for timeout
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}

	// check for DNS resolution failure
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return "dns_error"
	}

	// check for connection refused
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return "connection_error"
	}

	if strings.Contains(err.Error(), "unexpected status code") {
		return "http_error"
	}

	if strings.Contains(err.Error(), "decoding response") || strings.Contains(err.Error(), "parsing price") {
		return "parse_error"
	}

	return "unknown"
}
