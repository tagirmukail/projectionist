package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"projectionist/forms"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/session"
	"projectionist/utils"
)

func LoginApi(
	dbProvider provider.IDBProvider,
	sessHandler *session.SessionHandler,
	tokenSecretKey string,
) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		var respond = make(map[string]interface{})
		var form = forms.LoginForm{}
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			respond = utils.Message(false, "Invaid parameters")
			resp.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(resp, respond)
			return
		}

		err = form.Validate()
		if err != nil {
			respond = utils.Message(false, err.Error())
			resp.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(resp, respond)
			return
		}

		var user = models.User{}

		if err = dbProvider.GetByName(&user, form.Username); err != nil {
			respond = utils.Message(false, "User not exist")
			log.Printf("LoginApi() error: %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(resp, respond)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			respond = utils.Message(false, "Not authorized")
			log.Printf("LoginApi() bcrypt.CompareHashAndPassword() error: %v", err)
			resp.WriteHeader(http.StatusUnauthorized)
			utils.JsonRespond(resp, respond)
			return
		}

		var tokenM = &models.Token{UserId: uint64(user.ID)}
		var token = jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodHS256.Name), tokenM)
		tokenStr, err := token.SignedString([]byte(tokenSecretKey))
		if err != nil {
			respond = utils.Message(false, "Authorization failed")
			log.Printf("LoginApi() token.SignedString() error: %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(resp, respond)
			return
		}

		respond = utils.Message(true, "Login successful")
		user.Password = ""
		user.Token = tokenStr
		respond["user"] = user
		sessHandler.SetSession(user.Username, resp)
		utils.JsonRespond(resp, respond)

	})
}
