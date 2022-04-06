package server

import (
	"html/template"
	"net/http"
)

var templates = template.Must(
	template.ParseFiles(
		getTemplatePath("title"),
	))

func getTemplatePath(fileName string) string {
	templatesDir := "templates/"
	return templatesDir + fileName + ".html"
}

func renderTitleTemplate(w http.ResponseWriter, templateName string, titlePage *Title) {
	err := templates.ExecuteTemplate(w, templateName+".html", titlePage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
