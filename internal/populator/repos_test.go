package populator

import (
	"autobiography/internal/database"
	"autobiography/internal/models"
	"testing"
)

var testRepos = []models.GitHubRepo{
	{
		HtmlUrl:  "test.com/test",
		Language: "go",
		Name:     "test repo",
	},
	{
		HtmlUrl:  "test.com/test2",
		Language: "elixir",
		Name:     "test2 repo",
	},
}

func TestPopulateRepos(t *testing.T) {
	m, cleanup := models.SetupModels(t)
	candidate := models.Candidate{
		GivenName:  "test",
		FamilyName: "test",
	}
	_, err := m.Candidates.Insert(&candidate)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	err = PopulateRepos(m, 1, testRepos)
	if err != nil {
		t.Fatal(err)
	}

	checkReposEntries(t, m.ReposModel.DB, candidate.ID, testRepos)
	checkTechnologiesEntries(t, m.ReposModel.DB, testRepos)
}

func checkReposEntries(t *testing.T, db database.Db, candidateId int64, repos []models.GitHubRepo) {
	query := `
		SELECT html_url, technology_name, name
		FROM repos
	WHERE candidate_id = ? and name = ? and html_url = ? and technology_name = ?`

	var htmlUrl, technologyName, name string

	for _, repo := range repos {
		args := []interface{}{candidateId, repo.Name, repo.HtmlUrl, repo.Language}
		err := database.Get(db, query, args, []interface{}{&htmlUrl, &technologyName, &name})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func checkTechnologiesEntries(t *testing.T, db database.Db, repos []models.GitHubRepo) {
	query := `
		SELECT name
		FROM technologies
		WHERE name = ?`

	var name string

	for _, repo := range repos {
		args := []interface{}{repo.Language}
		err := database.Get(db, query, args, []interface{}{&name})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
