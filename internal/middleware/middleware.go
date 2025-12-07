// Package middleware defines common http middlewares.
package middleware

import (
	"net/http"
)

// Middleware take http.Handler,
// do something in between,
// then return http.Handler
type Middleware func(http.Handler) http.Handler

// CreateStack will struct the middlewares as:
//
//	middlewares := CreateStack(
//		middleware1,
//		middleware2,
//	)
//
// The sequence of middleware operated will be from top to bottom
func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
