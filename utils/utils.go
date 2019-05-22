package utils

import (
	"encoding/json"
	"net/http"
	"os"
)

// Message - represent struct of http response
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond - send for http client response
func Respond(w http.ResponseWriter, data map[string]interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
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
