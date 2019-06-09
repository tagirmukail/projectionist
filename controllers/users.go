package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

func NewUser(dbProvider provider.IDBProvider) http.HandlerFunc {
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

		if err, exist := dbProvider.IsExist(&user); exist && err == nil {
			w.WriteHeader(http.StatusForbidden)
			utils.Respond(w, utils.Message(false, "user exist"))
			return
		}

		err = dbProvider.Save(&user)
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

func GetUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func GetUserList(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func UpdateUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func DeleteUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}
