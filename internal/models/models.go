package models

import (
	"autobiography/internal/database"
)

type Models struct {
	Technologies          TechnologyModel
	Candidates            CandidateModel
	Communications        CommunicationModel
	EducationHistories    EducationHistoryModel
	EmploymentHistories   EmploymentHistoryModel
	TechnologyExperiences TechnologyExperienceModel
	ReposModel            GitHubReposModel
}

func NewModels(db database.Db) Models {
	return Models{
		Technologies:          TechnologyModel{db},
		Candidates:            CandidateModel{db},
		Communications:        CommunicationModel{db},
		EducationHistories:    EducationHistoryModel{db},
		EmploymentHistories:   EmploymentHistoryModel{db},
		TechnologyExperiences: TechnologyExperienceModel{db},
		ReposModel:            GitHubReposModel{db},
	}
}
