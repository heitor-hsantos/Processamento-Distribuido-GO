package middleware

import "net/http"

func Idempotency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("Idempotency-Key") == "" {
				http.Error(w, "missing idempotency key", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
