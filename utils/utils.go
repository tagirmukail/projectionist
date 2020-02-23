package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"projectionist/consts"
)

// Message - represent struct of http response
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// JsonRespond - send for http client response
func JsonRespond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error: json respond:%v", err)
	}
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

func Pagination(page, count int) (start, end int) {
	if page == 0 {
		return 0, 0
	}

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

func CheckHealthStatusCode(status int, name string) error {
	if status == http.StatusOK {
		return nil
	}

	if status == http.StatusUnauthorized {
		return fmt.Errorf("service %s not authorized", name)
	}

	if status == http.StatusBadRequest {
		return fmt.Errorf("service %s request bad", name)
	}

	if status == http.StatusNotFound {
		return fmt.Errorf("may be service %s is dead", name)
	}

	if status >= 500 {
		return fmt.Errorf("service %s dead", name)
	}

	return fmt.Errorf("service %s something when wrong", name)
}
