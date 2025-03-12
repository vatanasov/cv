package models

import (
	"autobiography/internal/database"
	"time"
)

type EducationHistory struct {
	OrganizationName string `xml:"OrganizationName"`
	Degree           string `xml:"EducationDegree>DegreeName"`
	From             CVDate `xml:"AttendancePeriod>StartDate>FormattedDateTime"`
	To               CVDate `xml:"AttendancePeriod>EndDate>FormattedDateTime"`
}

type EducationHistoryModel struct {
	DB database.Db
}

func (m *EducationHistoryModel) Insert(Candidate Candidate) error {
	query := `
		INSERT INTO education_histories (candidate_id, organization_name, degree, from_date, to_date)
		VALUES (?, ?, ?, ?, ?)`

	for _, educationHistory := range Candidate.EducationHistory {
		args := []any{
			Candidate.ID,
			educationHistory.OrganizationName,
			educationHistory.Degree,
			time.Time(educationHistory.From),
			time.Time(educationHistory.To),
		}
		_, err := database.Insert(m.DB, query, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *EducationHistoryModel) GetAll(candidateId int64) ([]EducationHistory, error) {
	query := `
		SELECT organization_name, degree, from_date, to_date
		FROM education_histories
		WHERE candidate_id = ?`

	args := []any{candidateId}
	resultF := func(educationHistory *EducationHistory) []any {
		return []any{
			&educationHistory.OrganizationName,
			&educationHistory.Degree,
			&educationHistory.From,
			&educationHistory.To,
		}
	}

	return database.GetAll(m.DB, query, args, resultF)
}
