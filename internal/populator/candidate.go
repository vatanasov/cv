package populator

import (
	"autobiography/internal/models"
)

func PopulateCandidate(models models.Models, candidate *models.Candidate) error {
	technologies := NewSet[string, internalTechnology]()
	for _, employmentHistory := range candidate.EmploymentHistory {
		for _, role := range employmentHistory.Description.Roles {
			for _, technology := range role.Technologies {
				technologies.Add(technology)
			}
		}
	}

	err := insertTechnologies(models, technologies)
	if err != nil {
		return err
	}

	_, err = models.Candidates.Insert(candidate)
	if err != nil {
		return err
	}

	err = models.Communications.Insert(*candidate)
	if err != nil {
		return err
	}

	err = models.EducationHistories.Insert(*candidate)
	if err != nil {
		return err
	}

	err = models.EmploymentHistories.Insert(candidate)
	if err != nil {
		return err
	}

	return nil
}
