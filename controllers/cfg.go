package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
	"strings"
)

func NewCfg(provider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		names, ok := r.URL.Query()["name"]
		if !ok || len(names) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.NameIsEmptyResp))
			return
		}

		var configFileName = names[0]
		if configFileName == "" {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.NameIsEmptyResp))
			return
		}
		configFileName = fmt.Sprintf("%v.json", configFileName)

		var form = make(map[string]interface{})
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		var cfg = &models.Configuration{}

		countCfgs, err := provider.Count(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		cfg.ID = countCfgs + 1
		cfg.Name = fmt.Sprintf(models.FileIDNAMEPtrn, cfg.ID, configFileName)
		cfg.Config = form
		err = cfg.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.InputDataInvalidResp))
			return
		}

		err = provider.Save(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("File %s saved", configFileName)
		utils.JsonRespond(w, utils.Message(true, msg))
	})
}

func GetCfg(provider provider.IDBProvider) http.HandlerFunc {
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

		var cfg models.Model
		cfg, err = provider.GetByID(cfg, int64(id))
		if err != nil {
			if err == consts.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
				return
			}
			log.Printf("GetCfg() error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		var respond = utils.Message(true, "")
		respond["config"] = cfg
		utils.JsonRespond(w, respond)
	})
}

func GetCfgList(provider provider.IDBProvider) http.HandlerFunc {
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

		countAllCfgs, err := provider.Count(&models.Configuration{})
		if err != nil {
			log.Printf("GetCfgList() count all configs error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		if start > countAllCfgs {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			utils.JsonRespond(w, respond)
			return
		}

		if end > countAllCfgs {
			end = countAllCfgs
		}

		cfgModels, err := provider.Pagination(&models.Configuration{}, start, end)
		if err != nil {
			log.Printf("GetCfgList() pagination by configs error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		var respond = utils.Message(true, "")
		respond[consts.KEY_CONFIGS] = cfgModels
		utils.JsonRespond(w, respond)
	})
}

func UpdateCfg(provider provider.IDBProvider) http.HandlerFunc {
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

		var cfg = models.Configuration{}
		err = json.NewDecoder(r.Body).Decode(&cfg)
		if err != nil {
			log.Printf("UpdateCfg() error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		err, exist := provider.IsExistByName(&cfg)
		if !exist {
			w.WriteHeader(http.StatusNotFound)
			utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		err = provider.Update(&cfg, id)
		if err != nil {
			if err == consts.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.NotUpdatedResp))
			return
		}

		w.WriteHeader(http.StatusOK)
		var respond = utils.Message(true, "Config updated")
		respond["config"] = cfg
		utils.JsonRespond(w, respond)
		return
	})
}

func DeleteCfg(provider provider.IDBProvider) http.HandlerFunc {
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

		var cfg = &models.Configuration{}
		err = provider.Delete(cfg, id)
		if err != nil {
			if err == consts.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.NotDeletedResp))
			return
		}

		w.WriteHeader(http.StatusOK)
		var respond = utils.Message(true, "config deleted")
		utils.JsonRespond(w, respond)
		return
	})
}
