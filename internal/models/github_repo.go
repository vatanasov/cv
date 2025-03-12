package models

import (
	"autobiography/internal/database"
)

type GitHubRepo struct {
	HtmlUrl  string `json:"html_url"`
	Language string `json:"language"`
	Name     string `json:"name"`
}

type GitHubReposModel struct {
	DB database.Db
}

func (m *GitHubReposModel) Insert(candidateId int64, repo GitHubRepo) error {
	query := `
		INSERT INTO repos (candidate_id, html_url, technology_name, name)
		VALUES (?, ?, ?, ?)`

	args := []interface{}{candidateId, repo.HtmlUrl, repo.Language, repo.Name}

	_, err := database.Insert(m.DB, query, args)
	return err
}
