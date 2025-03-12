package models

import (
	"autobiography/internal/database"
	"encoding/json"
)

type TechnologyExperience struct {
	Name           string
	PrettifiedName string
	OrderPriority  int
	Projects       TechnologyExperienceProjects
	Repos          GitHubRepos
}

type TechnologyExperienceProject struct {
	ID      int64  `json:"id"`
	Project string `json:"project"`
}

type TechnologyExperienceProjects []TechnologyExperienceProject

func (t *TechnologyExperienceProjects) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal([]byte(value.(string)), t)
}

type GitHubRepos []GitHubRepo

func (r *GitHubRepos) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal([]byte(value.(string)), r)
}

type TechnologyExperienceModel struct {
	Db database.Db
}

func (m *TechnologyExperienceModel) Get(candidateId int64, Name string) (TechnologyExperience, error) {
	query := `
		with ehp as (
			select ehrt.technology_name as name, json_group_array(json_object('id', ehr.id, 'project', ehr.project)) as projects
			from employment_histories eh
			join employment_history_roles ehr on eh.id = ehr.employment_history_id
			join employment_history_roles_technologies ehrt on ehr.id = ehrt.employment_history_role_id
			where eh.candidate_id = ?
			group by ehrt.technology_name
		), r as (
			select 
				r.technology_name as name, 
				json_group_array(json_object('language', r.technology_name, 'html_url', r.html_url, 'name', r.name)) as repos
			from repos r
			where r.candidate_id = ?
			group by r.technology_name
		)
		select t.name, t.prettified_name, t.order_priority, ehp.projects, r.repos
		from technologies t
		left join ehp on ehp.name = t.name
		left join r on r.name = t.name
		where t.name = ?
			and (ehp.projects is not null or r.repos is not null)
		order by t.order_priority`

	args := []any{candidateId, candidateId, Name}
	var technologyExperience TechnologyExperience
	result := []any{
		&technologyExperience.Name,
		&technologyExperience.PrettifiedName,
		&technologyExperience.OrderPriority,
		&technologyExperience.Projects,
		&technologyExperience.Repos,
	}

	err := database.Get(m.Db, query, args, result)

	return technologyExperience, err
}
