package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/session"
	"projectionist/utils"
	"strings"
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
				utils.JsonRespond(w, utils.Message(false, "Login required"))
				return
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}

var AccessControllAllows = func(accessAddresses []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var addrs string
			var addrsCount = len(accessAddresses)
			switch {
			case addrsCount == 1:
				addrs = accessAddresses[0]
			case addrsCount > 1:
				addrs = strings.Join(accessAddresses, ";")
			default:
				w.WriteHeader(http.StatusNetworkAuthenticationRequired)
				utils.JsonRespond(w, utils.Message(false, "Network not accepted"))
				return
			}

			w.Header().Add("Access-Control-Allow-Origin", addrs)
			w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Cookie")
			next.ServeHTTP(w, r)
		})
	}
}

var JwtAuthentication = func(tokenSecretKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestPath = r.URL.Path
			if requestPath == consts.UrlApiLogin {
				next.ServeHTTP(w, r)
			}

			var response = make(map[string]interface{})

			var tokenHeader = r.Header.Get(consts.AuthorizationHeader)

			if tokenHeader == "" {
				response = utils.Message(false, "Missing authorization token")
				w.WriteHeader(http.StatusForbidden)
				utils.JsonRespond(w, response)
				return
			}

			var splitted = strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				response = utils.Message(false, "Invalid/Malformed authorization token")
				w.WriteHeader(http.StatusForbidden)
				utils.JsonRespond(w, response)
				return
			}

			var tokenPart = splitted[1]
			var tokenM = &models.Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tokenM, func(token *jwt.Token) (i interface{}, e error) {
				return []byte(tokenSecretKey), nil
			})
			if err != nil {
				response = utils.Message(false, "Malformed authentication token")
				w.WriteHeader(http.StatusForbidden)
				utils.JsonRespond(w, response)
				return
			}

			if !token.Valid {
				response = utils.Message(false, "Authentication token is not valid")
				w.WriteHeader(http.StatusForbidden)
				utils.JsonRespond(w, response)
				return
			}
		})
	}
}
