package main

import (
	"html/template"
	"net/http"
	"time"
)

var TemplateFuncs = template.FuncMap{
	// Time functions
	"now":        time.Now,
	"timeSince":  time.Since,
	"timeUntil":  time.Until,
	"formatDate": formatDate,
	"raw":        raw,
}

func formatDate(t time.Time) string {
	return t.Format("2006")
}

func raw(s string) template.HTML {
	return template.HTML(s)
}

func (app *application) newTemplateData(r *http.Request) map[string]any {
	data := map[string]any{
		"WHAT": "TODO",
	}

	return data
}
