package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	writerKey  contextKey = "httpWriter"
	requestKey contextKey = "httpRequest"
)

func HTTPContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), writerKey, w)
		ctx = context.WithValue(ctx, requestKey, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Writer(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(writerKey).(http.ResponseWriter)
	return w, ok
}

func Request(ctx context.Context) (*http.Request, bool) {
	r, ok := ctx.Value(requestKey).(*http.Request)
	return r, ok
}
