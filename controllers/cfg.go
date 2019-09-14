package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

func NewCfg(provider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		names, ok := r.URL.Query()["name"]
		if !ok || len(names) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "name is empty"))
			return
		}

		var configFileName = names[0]
		if configFileName == "" {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "name is empty"))
			return
		}
		configFileName = fmt.Sprintf("%v.json", configFileName)

		var form = make(map[string]interface{})
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "Bad input data"))
			return
		}

		var cfg = &models.Configuration{}

		countCfgs, err := provider.Count(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "Something when wrong"))
			return
		}

		cfg.ID = countCfgs + 1
		cfg.Name = fmt.Sprintf(models.FileIDNAMEPtrn, cfg.ID, configFileName)
		cfg.Config = form
		err = provider.Save(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "Something when wrong, file not saved"))
			return
		}

		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("File %s saved", configFileName)
		utils.JsonRespond(w, utils.Message(true, msg))
	})
}

func GetCfg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func GetCfgList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func UpdateCfg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}

func DeleteCfg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: implement
	})
}
