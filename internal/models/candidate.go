package models

import (
	"autobiography/internal/database"
	"encoding/xml"
	"errors"
)

var (
	ErrInvalidRolesList = errors.New("invalid roles list")
)

type Candidate struct {
	ID                int64
	GivenName         string              `xml:"CandidatePerson>PersonName>GivenName"`
	FamilyName        string              `xml:"CandidatePerson>PersonName>FamilyName"`
	Communications    []Communication     `xml:"CandidatePerson>Communication"`
	EmploymentHistory []EmploymentHistory `xml:"CandidateProfile>EmploymentHistory>EmployerHistory"`
	EducationHistory  []EducationHistory  `xml:"CandidateProfile>EducationHistory>EducationOrganizationAttendance"`
}

func LoadCandidate(m Models) (Candidate, error) {
	candidate, err := m.Candidates.Get(1)
	if err != nil {
		return candidate, err
	}

	communications, err := m.Communications.GetAll(candidate.ID)
	if err != nil {
		return candidate, err
	}
	candidate.Communications = communications

	educationHistories, err := m.EducationHistories.GetAll(candidate.ID)
	if err != nil {
		return candidate, err
	}
	candidate.EducationHistory = educationHistories

	employmentHistories, err := m.EmploymentHistories.GetAll(candidate.ID)
	if err != nil {
		return candidate, err
	}
	candidate.EmploymentHistory = employmentHistories

	return candidate, nil
}

func FromXML(xmlContents []byte) (Candidate, error) {
	candidate := Candidate{}
	err := xml.Unmarshal(xmlContents, &candidate)

	return candidate, err
}

type CandidateModel struct {
	DB database.Db
}

func (m *CandidateModel) Insert(candidate *Candidate) (int64, error) {
	query := `
		INSERT INTO candidates (given_name, family_name)
		VALUES (?, ?)`

	args := []any{candidate.GivenName, candidate.FamilyName}

	id, err := database.Insert(m.DB, query, args)
	candidate.ID = id

	return id, err
}

func (m *CandidateModel) Get(id int64) (Candidate, error) {
	query := `
		SELECT given_name, family_name
		FROM candidates
		WHERE id = ?`

	args := []any{id}

	var candidate Candidate
	candidate.ID = id

	err := database.Get(m.DB, query, args, []any{&candidate.GivenName, &candidate.FamilyName})

	return candidate, err
}
