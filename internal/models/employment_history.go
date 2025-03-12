package models

import (
	"autobiography/internal/database"
	"encoding/xml"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	"time"
)

type EmploymentHistory struct {
	ID               int64
	OrganizationName string      `xml:"OrganizationName"`
	Position         string      `xml:"PositionHistory>PositionTitle"`
	Description      Description `xml:"PositionHistory>Description"`
	From             CVDate      `xml:"PositionHistory>EmploymentPeriod>StartDate>FormattedDateTime"`
	To               CVDate      `xml:"PositionHistory>EmploymentPeriod>EndDate>FormattedDateTime"`
	Currently        bool        `xml:"PositionHistory>EmploymentPeriod>CurrentIndicator"`
}

type Description struct {
	Text  string
	Roles []Role
}
type Role struct {
	ID           int64
	Project      string
	Role         string
	Technologies []Technology
}

func (r Role) PrintTechnologyStack() string {
	var technologies []string
	for _, technology := range r.Technologies {
		technologies = append(technologies, technology.TextEnhancement)
	}
	return strings.Join(technologies, ", ")
}

type EmploymentHistoryModel struct {
	DB database.Db
}

func (x *Description) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x.Roles = []Role{}
	x.Text = ""

	var rawDescription string
	err := d.DecodeElement(&rawDescription, &start)
	if err != nil {
		return err
	}

	myReader := strings.NewReader(rawDescription)

	nodes, err := html.ParseFragment(
		myReader, &html.Node{
			Type:     html.ElementNode,
			Data:     "div",
			DataAtom: atom.Div,
		},
	)

	if err != nil {
		return err
	}

	finalDescription := strings.Builder{}

	for _, node := range nodes {
		if node.Data == "ol" {
			err := processRoles(node, x)
			if err == nil {
				continue
			}
		}
		err := html.Render(&finalDescription, node)
		if err != nil {
			return err
		}
	}

	x.Text = finalDescription.String()

	return nil
}

func (m *EmploymentHistoryModel) Insert(Candidate *Candidate) error {
	query := `
		INSERT INTO employment_histories (candidate_id, organization_name, position, text_description, from_date, to_date, current)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	for i, employmentHistory := range Candidate.EmploymentHistory {
		args := []any{
			Candidate.ID,
			employmentHistory.OrganizationName,
			employmentHistory.Position,
			employmentHistory.Description.Text,
			time.Time(employmentHistory.From),
			time.Time(employmentHistory.To),
			employmentHistory.Currently,
		}
		id, err := database.Insert(m.DB, query, args)
		if err != nil {
			return err
		}

		Candidate.EmploymentHistory[i].ID = id

		err = insertRoles(m.DB, id, employmentHistory.Description.Roles)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *EmploymentHistoryModel) GetAll(candidateId int64) ([]EmploymentHistory, error) {
	query := `
		SELECT id, organization_name, position, text_description, from_date, to_date, current
		FROM employment_histories
		WHERE candidate_id = ?`

	args := []any{candidateId}
	resultF := func(employmentHistory *EmploymentHistory) []any {
		return []any{
			&employmentHistory.ID,
			&employmentHistory.OrganizationName,
			&employmentHistory.Position,
			&employmentHistory.Description.Text,
			&employmentHistory.From,
			&employmentHistory.To,
			&employmentHistory.Currently,
		}
	}

	employmentHistories, err := database.GetAll(m.DB, query, args, resultF)
	if err != nil {
		return employmentHistories, err
	}

	for i := range employmentHistories {
		employmentHistories[i].Description.Roles, err = getRoles(m.DB, employmentHistories[i].ID)
		if err != nil {
			return employmentHistories, err
		}
	}

	return employmentHistories, nil
}

func insertRoles(db database.Db, employmentHistoryID int64, roles []Role) error {
	query := `
		INSERT INTO employment_history_roles(employment_history_id, project, role)
		VALUES (?, ?, ?)`

	for _, role := range roles {
		args := []any{employmentHistoryID, role.Project, role.Role}
		id, err := database.Insert(db, query, args)
		if err != nil {
			return err
		}
		err = insertRoleTechnologies(db, id, role.Technologies)
		if err != nil {
			return err
		}
	}

	return nil
}

func getRoles(db database.Db, employmentHistoryID int64) ([]Role, error) {
	query := `
		SELECT id, project, role
		FROM employment_history_roles
		WHERE employment_history_id = ?`

	args := []any{employmentHistoryID}
	resultF := func(role *Role) []any {
		return []any{&role.ID, &role.Project, &role.Role}
	}

	roles, err := database.GetAll(db, query, args, resultF)
	if err != nil {
		return roles, err
	}

	for i := range roles {
		roles[i].Technologies, err = getRoleTechnologies(db, roles[i].ID)
		if err != nil {
			return roles, err
		}
	}

	return roles, nil
}

func getRoleTechnologies(db database.Db, roleId int64) ([]Technology, error) {
	query := `
		SELECT technology_name, prettified_name
		FROM employment_history_roles_technologies
		JOIN technologies ON employment_history_roles_technologies.technology_name = technologies.name
		WHERE employment_history_role_id = ?
		ORDER BY technologies.order_priority DESC`

	args := []any{roleId}
	resultF := func(technology *Technology) []any {
		return []any{&technology.Name, &technology.TextEnhancement}
	}

	return database.GetAll(db, query, args, resultF)
}

func insertRoleTechnologies(db database.Db, employmentHistoryRoleID int64, technologies []Technology) error {
	query := `
		INSERT INTO employment_history_roles_technologies(employment_history_role_id, technology_name)
		VALUES (?, ?)`

	for _, technology := range technologies {
		args := []any{employmentHistoryRoleID, technology.Name}
		_, err := database.Insert(db, query, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func processRoles(node *html.Node, d *Description) error {
	var roles []Role
	var text string
	childNode := node.FirstChild
	for childNode != nil {
		var role Role
		role.Project = collectText(childNode)
		childNode = childNode.NextSibling

		if childNode == nil {
			return ErrInvalidRolesList
		}
		text = collectText(childNode)
		if !strings.HasPrefix(text, "Role: ") {
			return ErrInvalidRolesList
		}
		role.Role = strings.TrimPrefix(text, "Role: ")
		childNode = childNode.NextSibling

		if childNode == nil {
			return ErrInvalidRolesList
		}
		text = collectText(childNode)
		if !strings.HasPrefix(text, "Stack: ") {
			return ErrInvalidRolesList
		}
		text = strings.TrimPrefix(text, "Stack: ")
		technologies := strings.Split(text, ", ")
		for _, technology := range technologies {
			technologyMeta := PrettyTechnologies.Get(technology)
			role.Technologies = append(
				role.Technologies,
				Technology{
					Name:            technologyMeta.Name,
					TextEnhancement: technologyMeta.Pretty,
					Order:           technologyMeta.Order,
				},
			)
		}
		childNode = childNode.NextSibling

		roles = append(roles, role)
	}
	d.Roles = roles
	return nil
}

func collectText(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}
	return ""
}
