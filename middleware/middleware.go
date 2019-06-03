package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"projectionist/consts"
	"projectionist/session"
	"projectionist/utils"
)

var LoginRequired = func(sessionHandler *session.SessionHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := sessionHandler.GetUserName(r)
			requestPath := r.URL.Path

			switch {
			case username == "" && requestPath != consts.UrlApiLogin:
				w.WriteHeader(http.StatusUnauthorized)
				utils.Respond(w, utils.Message(false, "Login required"))
				return
			case username != "" && requestPath == consts.UrlApiLogin:
				w.WriteHeader(http.StatusTeapot)
				utils.Respond(w, utils.Message(false, "Login complete"))
				return
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}
