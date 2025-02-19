package view

import (
	"app/util"
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/go-chi/render"
)

func Index(w http.ResponseWriter, r *http.Request, data any) {
	_, filename, _, _ := runtime.Caller(0)
	templateFile := filepath.Join(filepath.Dir(filename), "/template/index.html")
	renderHtml(w, r, templateFile, data)
}

func Error500(w http.ResponseWriter, r *http.Request, data any) {
	_, filename, _, _ := runtime.Caller(0)
	templateFile := filepath.Join(filepath.Dir(filename), "/template/error/500.html")
	renderHtml(w, r, templateFile, data)
}

func renderHtml(w http.ResponseWriter, r *http.Request, templateFile string, data any) {
	tmpl, err := template.ParseFiles(templateFile)
	util.CheckError(err, "Fail to get template")
	byteWriter := &bytes.Buffer{}
	tmpl.Execute(byteWriter, data)
	render.HTML(w, r, byteWriter.String())
}

func ServerError(w http.ResponseWriter, r *http.Request) {
	Error500(w, r, nil)
}
