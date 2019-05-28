package controllers

import (
	"html/template"
	"net/http"
)

func Login(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "base", nil)
	}
}
