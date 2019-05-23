package middleware

import (
	"net/http"
	"strings"
)

var Auth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auths := []string{"signin"}
		requestPath := r.URL.Path

		for _, value := range auths {
			if strings.Contains(requestPath, value) {
				next.ServeHTTP(w, r)
				return
			}
		}

		//TODO add check login
		//ctx := context.WithValue(r.Context(), "user", id)
		//r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
