package main

import (
	"autobiography/internal/models"
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	candidate, err := models.LoadCandidate(app.models)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	technologies, err := app.models.Technologies.GetAll()
	if err != nil || len(technologies) == 0 {
		app.serverError(w, r, err)
		return
	}

	candidateId := app.contextGetCandidateId(r)

	technologyExperience, err := app.models.TechnologyExperiences.Get(candidateId, technologies[0].Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data["Candidate"] = candidate
	data["Technologies"] = technologies
	data["TechnologyExperience"] = technologyExperience

	err = Page(w, http.StatusOK, data, "index.tmpl.html", "partials/experience.tmpl.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) technologyHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		app.serverError(w, r, fmt.Errorf("name is required"))
		return
	}

	data := app.newTemplateData(r)

	technologies, err := app.models.Technologies.GetAll()
	if err != nil || len(technologies) == 0 {
		app.serverError(w, r, err)
		return
	}

	candidateId := app.contextGetCandidateId(r)

	technologyExperience, err := app.models.TechnologyExperiences.Get(candidateId, name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data["Technologies"] = technologies
	data["TechnologyExperience"] = technologyExperience

	err = NamedTemplateWithHeaders(w, http.StatusOK, data, nil, "experience", "partials/experience.tmpl.html")
	if err != nil {
		app.serverError(w, r, err)
	}
}
