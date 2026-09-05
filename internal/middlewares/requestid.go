package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	requestIdKey ctxKey = iota
)

const (
	requestIdHeader = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIdHeader)

		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Add(requestIdHeader, id)
		ctx := context.WithValue(r.Context(), requestIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestIdFromContext(ctx context.Context) string {
	requestId := ctx.Value(requestIdKey).(string)
	return requestId
}
