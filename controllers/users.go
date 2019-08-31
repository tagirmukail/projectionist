package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
	"strconv"
)

//TODO: add tests for all user controllers

func NewUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user = models.User{}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Printf("new user decode request body error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "Bad input fields"))
			return
		}

		err = user.Validate()
		if err != nil {
			log.Printf("new user validate error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, err.Error()))
			return
		}

		err, exist := dbProvider.IsExist(&user)
		if exist && err == nil {
			w.WriteHeader(http.StatusForbidden)
			utils.JsonRespond(w, utils.Message(false, "User exist"))
			return
		}

		if err != nil {
			log.Printf("new user create error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "User not created"))
			return
		}

		err = dbProvider.Save(&user)
		if err != nil {
			log.Printf("new user save error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "User not saved"))
			return
		}

		respond := utils.Message(true, "New user created")
		respond["userID"] = user.ID

		utils.JsonRespond(w, respond)
	}
}

func GetUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		idStr, ok := params["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "id is empty"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "id is not number"))
			return
		}

		var user = models.User{}

		err = dbProvider.GetByID(&user, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(w, utils.Message(false, "user not exist"))
				return
			}
			log.Printf("GetUser() error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "server not "))
			return
		}

		user.Password = ""
		var respond = utils.Message(true, "")
		respond["user"] = user
		utils.JsonRespond(w, respond)
	})
}

func GetUserList(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result []models.User

		pageStr := r.URL.Query().Get(consts.PAGE_PARAM)
		countStr := r.URL.Query().Get(consts.COUNT_PARAM)
		if pageStr == "" || countStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "page and count required"))
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "page must be a number"))
			return
		}

		if page <= 0 {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			respond[consts.KEY_USERS] = result
			utils.JsonRespond(w, respond)
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "count must be a number"))
			return
		}

		start, end := utils.Pagination(page, count)

		countAllUsers, err := dbProvider.Count(&models.User{})
		if err != nil {
			log.Printf("GetUserList() count all users error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "unexpected error"))
			return
		}

		if start > countAllUsers {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			respond[consts.KEY_USERS] = result
			utils.JsonRespond(w, respond)
			return
		}

		if end > countAllUsers {
			end = countAllUsers
		}

		userModels, err := dbProvider.Pagination(&models.User{}, start, end)
		if err != nil {
			log.Printf("GetUserList() pagination by users error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "unexpected error"))
			return
		}

		var respond = utils.Message(true, "")
		respond[consts.KEY_USERS] = userModels
		utils.JsonRespond(w, respond)
	})
}

func UpdateUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		idStr, ok := params["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "id is empty"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "id is not number"))
			return
		}

		var user = models.User{}

		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "Bad input fields"))
			return
		}

		if err, exist := dbProvider.IsExist(&user); !exist || err != nil {
			w.WriteHeader(http.StatusForbidden)
			utils.JsonRespond(w, utils.Message(false, "user not exist"))
			return
		}

		err = dbProvider.Update(&user, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "user not updated"))
			return
		}
	})
}

func DeleteUser(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params = mux.Vars(r)
		idStr, ok := params["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "id is empty"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "id is not number"))
			return
		}

		var user = &models.User{}

		err = dbProvider.Delete(user, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "user not deleted"))
			return
		}
	})
}
