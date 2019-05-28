package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func ServicesIndex(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			resp["message"] = "Login failed"
		}

		t.ExecuteTemplate(w, "base", resp)
	}
}
