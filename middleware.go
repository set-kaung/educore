package main

import (
	"log/slog"
	"net/http"
)

type RouteChainer struct {
	routes []func(http.Handler) http.Handler
}

func NewRouteChainer(initial ...func(http.Handler) http.Handler) *RouteChainer {
	return &RouteChainer{routes: initial}
}

func (r *RouteChainer) Chain(next http.Handler) http.Handler {
	h := next
	for i := len(r.routes) - 1; i >= 0; i-- {
		h = r.routes[i](h)
	}
	return h
}

func (r *RouteChainer) Append(appendingRoutes ...func(http.Handler) http.Handler) *RouteChainer {
	newRoutes := make([]func(http.Handler) http.Handler, 0, len(r.routes)+len(appendingRoutes))
	newRoutes = append(newRoutes, r.routes...)
	newRoutes = append(newRoutes, appendingRoutes...)
	return &RouteChainer{routes: newRoutes}
}

func RequestLogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			slog.Info("request", "method", r.Method, "route", r.Pattern)
		}
		next.ServeHTTP(w, r)
	})
}
