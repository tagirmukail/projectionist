package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"projectionist/consts"
	"strconv"
	"strings"
)

// Message - represent struct of http response
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// JsonRespond - send for http client response
func JsonRespond(w http.ResponseWriter, data map[string]interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

// GetIDFromReq get id parameter from request
func GetIDFromReq(r *http.Request) (int, error) {
	var params = mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		return 0, fmt.Errorf(strings.ToLower(consts.IdIsEmptyResp))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf(strings.ToLower(consts.IdIsNotNumberResp))
	}

	return id, nil
}

// GetPageAndCountFromReq get page and count parameter from request
func GetPageAndCountFromReq(r *http.Request) (int, int, error) {
	pageStr := r.URL.Query().Get(consts.PAGE_PARAM)
	countStr := r.URL.Query().Get(consts.COUNT_PARAM)
	if pageStr == "" || countStr == "" {
		return 0, 0, fmt.Errorf(strings.ToLower(consts.PageAndCountRequiredResp))
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, fmt.Errorf(strings.ToLower(consts.PageMustNumberResp))
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, 0, fmt.Errorf(strings.ToLower(consts.CountMustNumberResp))
	}

	return page, count, nil
}

// CreateDir create dir
func CreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
	}

	return err
}

// SaveJsonFile save data in json file
func SaveJsonFile(path string, form map[string]interface{}) error {
	data, err := json.Marshal(form)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func Pagination(page, count int) (start, end int) {
	start = (page - 1) * count
	end = start + count
	return start, end
}

func GetFileName(path string) string {
	var parths = strings.Split(path, "/")
	var lastElemIndx = len(parths) - 1
	var lastElem = parths[lastElemIndx]
	return lastElem
}
