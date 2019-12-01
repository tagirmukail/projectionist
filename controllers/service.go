package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

func NewService(dbProvider provider.IDBProvider, syncShan chan string) http.HandlerFunc {
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

		syncShan <- service.Name

		respond := utils.Message(true, "New service created")
		respond["serviceID"] = service.ID

		utils.JsonRespond(w, respond)
	})
}

func GetService(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id, err = utils.GetIDFromReq(r)
		if err != nil {
			if err.Error() == strings.ToLower(consts.IdIsEmptyResp) {
				w.WriteHeader(http.StatusBadRequest)
				utils.JsonRespond(w, utils.Message(false, consts.IdIsEmptyResp))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.IdIsNotNumberResp))
			return
		}

		serviceModel, err := dbProvider.GetByID(&models.Service{}, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("service with id %v not exist", id)
				w.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
				return
			}
			log.Printf("GetService() error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		service, ok := serviceModel.(*models.Service)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
			return
		}

		service.Emails, err = service.GetEmails()
		if err != nil {
			log.Printf("error: get service: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		var respond = utils.Message(true, "")
		respond["service"] = service
		utils.JsonRespond(w, respond)
	})
}

func GetServiceList(dbProvider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, count, err := utils.GetPageAndCountFromReq(r)
		if err != nil {
			if err.Error() == strings.ToLower(consts.PageAndCountRequiredResp) {
				w.WriteHeader(http.StatusBadRequest)
				utils.JsonRespond(w, utils.Message(false, consts.PageAndCountRequiredResp))
				return
			}

			if err.Error() == strings.ToLower(consts.PageMustNumberResp) {
				w.WriteHeader(http.StatusBadRequest)
				utils.JsonRespond(w, utils.Message(false, consts.PageMustNumberResp))
				return
			}

			if err.Error() == strings.ToLower(consts.CountMustNumberResp) {
				w.WriteHeader(http.StatusBadRequest)
				utils.JsonRespond(w, utils.Message(false, consts.CountMustNumberResp))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		if page <= 0 {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			utils.JsonRespond(w, respond)
			return
		}

		start, end := utils.Pagination(page, count)

		countAllServices, err := dbProvider.Count(&models.Service{})
		if err != nil {
			log.Printf("dbProvider.Count() count all services error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		if start > countAllServices {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			utils.JsonRespond(w, respond)
			return
		}

		if end > countAllServices {
			end = countAllServices
		}

		serviceModels, err := dbProvider.Pagination(&models.Service{}, start, end)
		if err != nil {
			log.Printf("dbProvider.Pagination() pagination by services error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		var respond = utils.Message(true, "")
		respond[consts.KEY_SERVICES] = serviceModels
		utils.JsonRespond(w, respond)
	})
}

func UpdateService(dbProvider provider.IDBProvider, syncShan chan string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id, err = utils.GetIDFromReq(r)
		if err != nil {
			if err.Error() == strings.ToLower(consts.IdIsEmptyResp) {
				w.WriteHeader(http.StatusBadRequest)
				utils.JsonRespond(w, utils.Message(false, consts.IdIsEmptyResp))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.IdIsNotNumberResp))
			return
		}

		var service = models.Service{ID: id}
		err = json.NewDecoder(r.Body).Decode(&service)
		if err != nil {
			log.Printf("decode request body error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		err = dbProvider.Update(&service, id)
		if err != nil {
			log.Printf("dbProvider.Update update service error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.NotUpdatedResp))
			return
		}

		iService, err := dbProvider.GetByID(&models.Service{}, int64(id))
		if err != nil {
			log.Printf("dbProvider.GetByID error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		syncShan <- iService.GetName()

		var respond = utils.Message(true, "service updated")
		respond["service"] = service
		utils.JsonRespond(w, respond)
	})
}

func DeleteService(dbProvider provider.IDBProvider, syncChan chan string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id, err = utils.GetIDFromReq(r)
		if err != nil {
			if err.Error() == strings.ToLower(consts.IdIsEmptyResp) {
				w.WriteHeader(http.StatusBadRequest)
				utils.JsonRespond(w, utils.Message(false, consts.IdIsEmptyResp))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.IdIsNotNumberResp))
			return
		}

		iService, err := dbProvider.GetByID(&models.Service{}, int64(id))
		if err != nil {
			log.Printf("dbProvider.GetByID error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		err = dbProvider.Delete(&models.Service{}, id)
		if err != nil {
			log.Printf("projectionist-api: dbProvider.Delete error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.NotDeletedResp))
			return
		}

		syncChan <- iService.GetName()

		var respond = utils.Message(true, "service deleted")
		utils.JsonRespond(w, respond)
	})
}
