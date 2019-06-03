package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"projectionist/models"
	"projectionist/utils"
)

func NewUser(sqlDB *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user = models.User{}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			utils.Respond(w, utils.Message(false, "Bad input fields"))
			return
		}

		err = user.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.Respond(w, utils.Message(false, err.Error()))
			return
		}

		if err, exist := user.IsExist(sqlDB); exist && err == nil {
			w.WriteHeader(http.StatusForbidden)
			utils.Respond(w, utils.Message(false, "user exist"))
			return
		}

		err = user.Save(sqlDB)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.Respond(w, utils.Message(false, "user not saved"))
			return
		}

		respond := utils.Message(true, "New user created")
		respond["userID"] = user.ID

		utils.Respond(w, respond)
	})
}
