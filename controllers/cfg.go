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
			utils.Respond(w, utils.Message(false, "name is empty"))
			return
		}

		var configFileName = names[0]
		configFileName = configFileName + ".json"

		var form = make(map[string]interface{})
		var err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.Respond(w, utils.Message(false, "Bad input data"))
			return
		}
		var savePath = consts.PathSaveCfgs + "/" + configFileName
		err = utils.SaveJsonFile(savePath, form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Respond(w, utils.Message(false, "File not saved"))
			return
		}

		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("File %s saved", configFileName)
		utils.Respond(w, utils.Message(true, msg))
	})
}
