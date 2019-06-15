package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"projectionist/consts"
	"projectionist/session"
	"projectionist/utils"
)

//LoginRequired middleware login
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
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}

var AccessControllAllows = func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Cookie")
			next.ServeHTTP(w, r)
		})
	}
}
