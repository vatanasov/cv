package models

import (
	"autobiography/internal/database"
	"testing"
)

func TestTechnologyModel_Insert(t *testing.T) {
	t.Run(
		"Insert a single struct", func(t *testing.T) {
			models, cleanup := SetupModels(t)
			defer cleanup()

			technology := Technology{
				Name:            "Go",
				TextEnhancement: "Go",
				Order:           1,
			}

			err := models.Technologies.Insert(technology)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !checkTechnologyExistence(models.Technologies.DB, "Go") {
				t.Fatal("Technology not found in database")
			}
		},
	)

	t.Run(
		"Ignore insert when struct already exists", func(t *testing.T) {
			models, cleanup := SetupModels(t)
			defer cleanup()

			technology := Technology{
				Name:            "Go",
				TextEnhancement: "Go",
				Order:           1,
			}

			err := models.Technologies.Insert(technology)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			err = models.Technologies.Insert(technology)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !checkTechnologyExistence(
				models.Technologies.DB,
				"Go",
			) || numberOfTechnologies(models.Technologies.DB) != 1 {
				t.Fatal("Technology not found in database")
			}
		},
	)
}

func TestTechnologyModel_GetAll(t *testing.T) {
	t.Run(
		"Get all technologies", func(t *testing.T) {
			models, cleanup := SetupModels(t)
			defer cleanup()

			result, err := models.Technologies.GetAll()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != 0 {
				t.Fatalf("expected 0 technologies, got %d", len(result))
			}

			technologies := []Technology{
				{
					Name:            "Go",
					TextEnhancement: "Go",
					Order:           1,
				},
				{
					Name:            "Python",
					TextEnhancement: "Python",
					Order:           2,
				},
			}

			for _, technology := range technologies {
				err := models.Technologies.Insert(technology)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			result, err = models.Technologies.GetAll()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != 2 {
				t.Fatalf("expected 2 technologies, got %d", len(result))
			}
		},
	)
}

func checkTechnologyExistence(db database.Db, name string) bool {
	query := `
		SELECT name
		FROM technologies
		WHERE name = ?`

	args := []any{name}
	var result string

	err := database.Get(db, query, args, []any{&result})
	return err == nil && result != ""
}

func numberOfTechnologies(db database.Db) int {
	query := `
		SELECT count(*)
		FROM technologies`

	args := []any{}
	var result int

	err := database.Get(db, query, args, []any{&result})
	if err != nil {
		panic(err)
	}

	return result
}
