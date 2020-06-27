package vk

import (
	"net/http"

	"github.com/suborbital/vektor/vlog"
)

// Middleware represents a handler that runs on a request before reaching its handler
type Middleware func(*http.Request, *Ctx) error

// ContentTypeMiddleware allows the content-type to be set
func ContentTypeMiddleware(contentType string) Middleware {
	return func(r *http.Request, ctx *Ctx) error {
		ctx.Headers.Set(contentTypeHeaderKey, contentType)

		return nil
	}
}

// CORSMiddleware enables CORS with the given domain for a route
// pass "*" to allow all domains, or empty string to allow none
func CORSMiddleware(domain string) Middleware {
	return func(r *http.Request, ctx *Ctx) error {
		enableCors(ctx, domain)

		return nil
	}
}

// CORSHandler enables CORS for a route
// pass "*" to allow all domains, or empty string to allow none
func CORSHandler(domain string) HandlerFunc {
	return func(r *http.Request, ctx *Ctx) (interface{}, error) {
		enableCors(ctx, domain)

		return nil, nil
	}
}

func enableCors(ctx *Ctx, domain string) {
	if domain != "" {
		ctx.Headers.Set("Access-Control-Allow-Origin", domain)
		ctx.Headers.Set("X-Requested-With", "XMLHttpRequest")
		ctx.Headers.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, cache-control")
	}
}

func loggerMiddleware(logger vlog.Logger) Middleware {
	return func(r *http.Request, ctx *Ctx) error {
		logger.Info(r.Method, r.URL.String())

		return nil
	}
}

// generate a HandlerFunc that passes the request through a set of Middleware first
func handlerWithMiddleware(inner HandlerFunc, middleware []Middleware) HandlerFunc {
	return func(r *http.Request, ctx *Ctx) (interface{}, error) {
		for _, m := range middleware {
			if err := m(r, ctx); err != nil {
				return nil, err
			}
		}

		return inner(r, ctx)
	}
}
