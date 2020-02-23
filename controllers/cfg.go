package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/grpclog"

	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

func NewCfg(provider provider.IDBProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cfgForm = &models.Configuration{}
		var err = json.NewDecoder(r.Body).Decode(cfgForm)
		if err != nil {
			grpclog.Errorf("decode form error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		err = cfgForm.Validate()
		if err != nil {
			grpclog.Errorf("validate form error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.InputDataInvalidResp))
			return
		}

		err = provider.Save(cfgForm)
		if err != nil {
			grpclog.Errorf("save form error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("File %s saved", cfgForm.Name)
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
			grpclog.Errorf("utils.GetIDFromReq() error: %v", err)
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
			grpclog.Errorf("provider.GetByID() error: %v", err)
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

			grpclog.Errorf("utils.GetPageAndCountFromReq error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		if page <= 0 {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			respond[consts.KEY_CONFIGS] = make([]models.Model, 0)
			utils.JsonRespond(w, respond)
			return
		}

		start, end := utils.Pagination(page, count)

		countAllCfgs, err := provider.Count(&models.Configuration{})
		if err != nil {
			grpclog.Errorf("GetCfgList() count all configs error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		if start > countAllCfgs {
			w.WriteHeader(http.StatusOK)
			respond := utils.Message(true, "")
			respond[consts.KEY_CONFIGS] = make([]models.Model, 0)
			utils.JsonRespond(w, respond)
			return
		}

		if end > countAllCfgs {
			end = countAllCfgs
		}

		cfgModels, err := provider.Pagination(&models.Configuration{}, start, end)
		if err != nil {
			grpclog.Errorf("GetCfgList() pagination by configs error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
			return
		}

		if cfgModels == nil {
			cfgModels = make([]models.Model, 0)
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
			grpclog.Errorf("utils.GetIDFromReq() error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.IdIsNotNumberResp))
			return
		}

		var cfg = models.Configuration{}
		err = json.NewDecoder(r.Body).Decode(&cfg)
		if err != nil {
			grpclog.Errorf("decode request body error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		err = cfg.Validate()
		if err != nil {
			grpclog.Errorf("cfg.Validate() error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.JsonRespond(w, utils.Message(false, consts.BadInputDataResp))
			return
		}

		err = provider.Update(&cfg, id)
		if err != nil {
			if err == consts.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				utils.JsonRespond(w, utils.Message(false, consts.NotExistResp))
				return
			}

			grpclog.Errorf("provider.Update(id:%d) error: %v", id, err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.JsonRespond(w, utils.Message(false, consts.SmtWhenWrongResp))
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
			grpclog.Errorf("utils.GetIDFromReq() error: %v", err)
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

			grpclog.Errorf("provider.Delete(id:%d) error: %v", id, err)
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
