package extractor

import (
	"autobiography/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
)

type Token string

type GitHubPin struct {
	Data struct {
		Viewer struct {
			PinnedItems struct {
				Nodes []struct {
					Name            string `json:"name"`
					Url             string `json:"url"`
					PrimaryLanguage struct {
						Name string `json:"name"`
					} `json:"primaryLanguage"`
				} `json:"nodes"`
			} `json:"pinnedItems"`
		} `json:"viewer"`
	} `json:"data"`
}

type RepoApiClient struct {
	BaseUrl string
}

var GitHubRepoApiClient = &RepoApiClient{
	BaseUrl: "https://api.github.com/graphql",
}

func (r *RepoApiClient) ExtractFromGitHub(token Token) ([]models.GitHubRepo, error) {
	query := map[string]string{
		"query": `query {
		  viewer {
			pinnedItems(first: 100, types: REPOSITORY) {
			  nodes {
				... on Repository {
				  name
				  url
				  primaryLanguage {
					name
				  }
				}
			  }
			}
		  }
		}`,
	}
	reqBody, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, r.BaseUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+string(token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var gitHubPins GitHubPin
	err = json.NewDecoder(res.Body).Decode(&gitHubPins)
	if err != nil {
		return nil, err
	}

	repos := make([]models.GitHubRepo, len(gitHubPins.Data.Viewer.PinnedItems.Nodes))
	for i, pin := range gitHubPins.Data.Viewer.PinnedItems.Nodes {
		var repo models.GitHubRepo
		repo.HtmlUrl = pin.Url
		repo.Language = pin.PrimaryLanguage.Name
		repo.Name = pin.Name
		repos[i] = repo
	}

	return repos, nil
}
