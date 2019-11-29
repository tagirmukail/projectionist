package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

func NewService(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var service = models.Service{}

		err := json.NewDecoder(r.Body).Decode(&service)
		if err != nil {
			log.Printf("new service decode request body error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		err = service.Validate()
		if err != nil {
			log.Printf("new service validate error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.InputDataInvalidResp))
			return
		}

		err, exist := dbProvider.IsExistByName(&service)
		if exist && err == nil {
			w.WriteHeader(http.StatusForbidden)
			utils.JsonRespond(w, utils.Message(false, "A service with the same name already exists."))
			return
		}

		if err != nil {
			log.Printf("new service create error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		err = dbProvider.Save(&service)
		if err != nil {
			log.Printf("new service save error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.NotSavedResp))
			return
		}

		respond := utils.Message(true, "New service created")
		respond["serviceID"] = service.ID

		utils.JsonRespond(w, respond)
	})
}

func GetService(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func GetServiceList(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func UpdateService(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func DeleteService(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}
