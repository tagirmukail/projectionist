package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"projectionist/consts"
	"projectionist/models"
	"projectionist/utils"
)

var JwtAuthentication = func(tokenSecretKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestPath = r.URL.Path
			if requestPath == consts.UrlApiLoginV1 || strings.Contains(requestPath, "swagger") {
				next.ServeHTTP(w, r)
				return
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

			next.ServeHTTP(w, r)
		})
	}
}
