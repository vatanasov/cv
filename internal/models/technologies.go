package models

import (
	"autobiography/internal/database"
	"errors"
)

type Technology struct {
	Name            string
	TextEnhancement string
	Order           int
}

func (t Technology) Key() string {
	return t.Name
}

func (t Technology) Value() Technology {
	return t
}

func (t Technology) String() string {
	return t.TextEnhancement
}

type TechnologyModel struct {
	DB database.Db
}

func (m *TechnologyModel) Insert(technology Technology) error {
	query := `
		INSERT INTO technologies (name, prettified_name, order_priority)
		VALUES (?, ?, ?)
		ON CONFLICT (name) DO NOTHING`
	args := []any{technology.Name, technology.TextEnhancement, technology.Order}

	_, err := database.Insert(m.DB, query, args)

	if err != nil {
		switch {
		case errors.Is(err, database.ErrNoRowsAffected):
			return nil
		default:
			return err
		}
	}

	return nil
}

func (m *TechnologyModel) GetAll() ([]Technology, error) {
	query := `
		SELECT name, prettified_name, order_priority
		FROM technologies
		ORDER BY order_priority DESC`

	args := []any{}
	resultF := func(technology *Technology) []any {
		return []any{&technology.Name, &technology.TextEnhancement, &technology.Order}
	}

	return database.GetAll(m.DB, query, args, resultF)
}
