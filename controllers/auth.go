package controllers

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"projectionist/forms"
	"projectionist/models"
	"projectionist/utils"
)

func LoginApi(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		var form = forms.LoginForm{}
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			var respond = utils.Message(false, "Invaid parameters")
			resp.WriteHeader(http.StatusBadRequest)
			utils.Respond(resp, respond)
			return
		}

		err = form.Validate()
		if err != nil {
			var respond = utils.Message(false, err.Error())
			resp.WriteHeader(http.StatusBadRequest)
			utils.Respond(resp, respond)
			return
		}

		var user = models.User{}

		if err = user.GetByName(db, form.Username); err != nil {
			var respond = utils.Message(false, "User not exist")
			log.Printf("LoginApi() error: %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			utils.Respond(resp, respond)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			var respond = utils.Message(false, "Not authorized")
			log.Printf("LoginApi() error: %v", err)
			resp.WriteHeader(http.StatusUnauthorized)
			utils.Respond(resp, respond)
			return
		}

		var respond = utils.Message(true, "Login successful")
		user.Password = ""
		respond["user"] = user
		utils.Respond(resp, respond)

	})
}
