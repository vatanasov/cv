package populator

import (
	"autobiography/internal/database"
	"autobiography/internal/models"
	"testing"
	"time"
)

var technologies = []models.Technology{
	{
		Name:            "Go",
		TextEnhancement: "Go pretty",
		Order:           9,
	},
}

var testCandidate = models.Candidate{
	GivenName:  "John",
	FamilyName: "Doe",
	Communications: []models.Communication{
		{
			ChannelCode: "email",
			URI:         "test@test.com",
		},
	},
	EmploymentHistory: []models.EmploymentHistory{
		{
			OrganizationName: "Test organization",
			Position:         "Test position",
			Description: models.Description{
				Text: "Test text for description",
				Roles: []models.Role{
					{
						Project:      "Test project",
						Role:         "Test role",
						Technologies: technologies,
					},
				},
			},
			From:      models.CVDate(time.Now()),
			To:        models.CVDate(time.Now()),
			Currently: false,
		},
	},
	EducationHistory: []models.EducationHistory{
		{
			OrganizationName: "University of Test",
			Degree:           "Test degree",
			From:             models.CVDate(time.Now()),
			To:               models.CVDate(time.Now()),
		},
	},
}

func TestPopulateCandidate(t *testing.T) {
	m, cleanup := models.SetupModels(t)
	defer cleanup()

	err := PopulateCandidate(m, &testCandidate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkCandidateEntry(t, m.Candidates.DB, testCandidate)
	checkCommunicationEntries(t, m.Communications.DB, testCandidate)
	checkEmploymentHistoryEntries(t, m.EmploymentHistories.DB, testCandidate)
	checkEducationHistoryEntries(t, m.EducationHistories.DB, testCandidate)
	checkRoleEntries(t, m.EmploymentHistories.DB, testCandidate)
	checkTechnologies(t, m.Technologies.DB)
}

func checkCandidateEntry(t *testing.T, db database.Db, candidate models.Candidate) {
	query := `
		SELECT given_name, family_name
		FROM candidates
		WHERE id = ?`

	args := []interface{}{candidate.ID}

	var givenName, familyName string
	err := database.Get(db, query, args, []interface{}{&givenName, &familyName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if givenName != candidate.GivenName {
		t.Fatalf("unexpected given name: %v", givenName)
	}

	if familyName != candidate.FamilyName {
		t.Fatalf("unexpected family name: %v", familyName)
	}
}

func checkCommunicationEntries(t *testing.T, db database.Db, candidate models.Candidate) {
	query := `
		SELECT channel_code, uri
		FROM communications
		WHERE candidate_id = ?`

	args := []interface{}{candidate.ID}

	var channelCode, uri string
	err := database.Get(db, query, args, []interface{}{&channelCode, &uri})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if channelCode != candidate.Communications[0].ChannelCode {
		t.Fatalf("unexpected channel code: %v", channelCode)
	}

	if uri != candidate.Communications[0].URI {
		t.Fatalf("unexpected uri: %v", uri)
	}
}

func checkEmploymentHistoryEntries(t *testing.T, db database.Db, candidate models.Candidate) {
	query := `
		SELECT organization_name, position, text_description, from_date, to_date, current
		FROM employment_histories
		WHERE candidate_id = ?`

	args := []interface{}{candidate.ID}

	var organizationName, position, textDescription string
	var fromDate, toDate time.Time
	var current bool
	err := database.Get(
		db,
		query,
		args,
		[]interface{}{&organizationName, &position, &textDescription, &fromDate, &toDate, &current},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if organizationName != candidate.EmploymentHistory[0].OrganizationName {
		t.Fatalf("unexpected organization name: %v", organizationName)
	}

	if position != candidate.EmploymentHistory[0].Position {
		t.Fatalf("unexpected position: %v", position)
	}

	if textDescription != candidate.EmploymentHistory[0].Description.Text {
		t.Fatalf("unexpected text description: %v", textDescription)
	}

	if !datesAreEqual(fromDate, time.Time(candidate.EmploymentHistory[0].From)) {
		t.Fatalf("want: %v, got: %v", time.Time(candidate.EmploymentHistory[0].From), fromDate)
	}

	if !datesAreEqual(toDate, time.Time(candidate.EmploymentHistory[0].To)) {
		t.Fatalf("unexpected to date: %v", toDate)
	}

	if current != candidate.EmploymentHistory[0].Currently {
		t.Fatalf("unexpected current: %v", current)
	}
}

func checkRoleEntries(t *testing.T, db database.Db, candidate models.Candidate) {
	query := `
		SELECT project, role
		FROM employment_history_roles
		WHERE employment_history_id = ?`

	args := []interface{}{candidate.EmploymentHistory[0].ID}

	var project, role string
	err := database.Get(db, query, args, []interface{}{&project, &role})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project != candidate.EmploymentHistory[0].Description.Roles[0].Project {
		t.Fatalf("unexpected project: %v", project)
	}

	if role != candidate.EmploymentHistory[0].Description.Roles[0].Role {
		t.Fatalf("unexpected role: %v", role)
	}
}

func checkTechnologies(t *testing.T, db database.Db) {
	query := `
		SELECT name, prettified_name, order_priority
		FROM technologies`

	args := []interface{}{}

	var name, textEnhancement string
	var order int
	err := database.Get(db, query, args, []interface{}{&name, &textEnhancement, &order})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name != technologies[0].Name {
		t.Fatalf("unexpected name: %v", name)
	}

	if textEnhancement != technologies[0].TextEnhancement {
		t.Fatalf("unexpected text enhancement: %v", textEnhancement)
	}

	if order != technologies[0].Order {
		t.Fatalf("unexpected order: %v", order)
	}
}

func checkEducationHistoryEntries(t *testing.T, db database.Db, candidate models.Candidate) {
	query := `
		SELECT organization_name, degree, from_date, to_date
		FROM education_histories
		WHERE candidate_id = ?`

	args := []interface{}{candidate.ID}

	var organizationName, degree string
	var fromDate, toDate time.Time
	err := database.Get(db, query, args, []interface{}{&organizationName, &degree, &fromDate, &toDate})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if organizationName != candidate.EducationHistory[0].OrganizationName {
		t.Fatalf("unexpected organization name: %v", organizationName)
	}

	if degree != candidate.EducationHistory[0].Degree {
		t.Fatalf("unexpected degree: %v", degree)
	}

	if !datesAreEqual(fromDate, time.Time(candidate.EducationHistory[0].From)) {
		t.Fatalf("unexpected from date: %v", fromDate)
	}

	if !datesAreEqual(toDate, time.Time(candidate.EducationHistory[0].To)) {
		t.Fatalf("unexpected to date: %v", toDate)
	}
}

func datesAreEqual(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
