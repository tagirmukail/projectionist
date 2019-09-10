package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectionist/consts"
	"projectionist/utils"
)

func NewCfg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		names, ok := r.URL.Query()["name"]
		if !ok || len(names) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "name is empty"))
			return
		}

		var configFileName = names[0]
		configFileName = fmt.Sprintf("%v.json", configFileName)

		var form = make(map[string]interface{})
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, "Bad input data"))
			return
		}
		var savePath = fmt.Sprintf("%s/%s", consts.PathSaveCfgs, configFileName)
		err = utils.SaveJsonFile(savePath, form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, "File not saved"))
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
