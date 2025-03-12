package populator

import (
	"autobiography/internal/models"
)

func PopulateRepos(m models.Models, candidateId int64, repos []models.GitHubRepo) error {
	technologies := NewSet[string, internalTechnology]()
	for i, repo := range repos {
		technologyMeta := models.PrettyTechnologies.Get(repo.Language)
		technology := internalTechnology{
			Name:            technologyMeta.Name,
			TextEnhancement: technologyMeta.Pretty,
			Order:           technologyMeta.Order,
		}

		repos[i].Language = technologyMeta.Name

		technologies.Add(technology)
	}

	err := insertTechnologies(m, technologies)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		err := m.ReposModel.Insert(candidateId, repo)
		if err != nil {
			return err
		}
	}
	return nil
}
