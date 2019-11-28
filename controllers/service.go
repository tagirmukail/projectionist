package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/utils"
)

func NewService() http.HandlerFunc {
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
		}

		//TODO: continue

	})
}

func GetService() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func GetServiceList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func UpdateService() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func DeleteService() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}
