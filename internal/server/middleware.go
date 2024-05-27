package server

import (
	"net/http"

	"github.com/juju/ratelimit"
)

func RateLimit(limit int, burst int, next http.Handler) http.Handler {
	bucket := ratelimit.NewBucketWithRate(float64(limit), int64(burst))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if bucket.TakeAvailable(1) == 0 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
