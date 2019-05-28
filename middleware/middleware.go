package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"projectionist/consts"
	"projectionist/session"
	"strings"
)

var LoginRequired = func(sessionHandler *session.SessionHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := sessionHandler.GetUserName(r)
			requestPath := r.URL.Path

			switch {
			case username == "" && !strings.Contains(requestPath, "login"):
				http.Redirect(w, r, "/login", 301)
				return
			case username != "" && strings.Contains(requestPath, "login"):
				http.Redirect(w, r, "/", 301)
				return
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
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
