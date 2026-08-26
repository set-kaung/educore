package web

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

func SameOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !safeMethod(r.Method) && !sameSite(r) {
			http.Error(w, "cross-site request blocked", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func sameSite(r *http.Request) bool {
	switch site := r.Header.Get("Sec-Fetch-Site"); site {
	case "":
	case "same-origin", "none":
	default:
		return false
	}

	origin := r.Header.Get("Origin")
	if origin == "" || origin == "null" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host

	return strings.EqualFold(u.Scheme, scheme) && strings.EqualFold(u.Hostname(), hostname(host)) && port(u.Port(), portFromHost(host))
}

func safeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func hostname(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

func portFromHost(host string) string {
	if _, p, err := net.SplitHostPort(host); err == nil {
		return p
	}
	return ""
}

func port(a, b string) bool {
	return a == b
}
