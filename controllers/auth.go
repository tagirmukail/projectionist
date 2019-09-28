package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"projectionist/consts"
	"projectionist/forms"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

func LoginApi(
	dbProvider provider.IDBProvider,
	tokenSecretKey string,
) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		var respond = make(map[string]interface{})
		var form = forms.LoginForm{}
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			respond = utils.Message(false, consts.InputDataInvalidResp)
			resp.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(resp, respond)
			return
		}

		err = form.Validate()
		if err != nil {
			respond = utils.Message(false, consts.InputDataInvalidResp)
			resp.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(resp, respond)
			return
		}

		var iUser models.Model
		var user = &models.User{}

		iUser, err = dbProvider.GetByName(user, form.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				respond = utils.Message(false, consts.NotExistResp)
				resp.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(resp, respond)
				return
			}

			respond = utils.Message(false, consts.SmtWhenWrongResp)
			log.Printf("LoginApi() error: %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(resp, respond)
			return
		}

		user, ok := iUser.(*models.User)
		if !ok {
			respond = utils.Message(false, consts.SmtWhenWrongResp)
			log.Printf("LoginApi() error: iUser is not User model")
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
		utils.JsonRespond(resp, respond)

	})
}
