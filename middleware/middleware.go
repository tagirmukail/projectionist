package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"projectionist/consts"
	"strings"
)

var Auth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		authPathParts := []string{"login"}
		for _, pathPart := range authPathParts {
			if strings.Contains(requestPath, pathPart) {

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

var FirstAuth = func(usersNotEmpty *bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "UsersNotEmpty", usersNotEmpty)
			r = r.WithContext(ctx)

			if !*usersNotEmpty {
				http.Redirect(w, r, consts.UrlNewUser, 301)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
