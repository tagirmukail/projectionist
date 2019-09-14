package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
	"strconv"
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
		var params = mux.Vars(r)
		idStr, ok := params["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.IdIsEmptyResp))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
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
		respond["cfg"] = cfg
		utils.JsonRespond(w, respond)
	})
}

func GetCfgList(provider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func UpdateCfg(provider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func DeleteCfg(provider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}
