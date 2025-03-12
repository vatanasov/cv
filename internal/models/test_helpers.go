package models

import (
	"autobiography/internal/database"
	"testing"
)

func SetupModels(t *testing.T) (Models, func()) {
	t.Helper()
	db, err := database.New("../testdata/testdb.sqlite")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	return NewModels(tx), func() {
		tx.Rollback()
		db.Close()
	}
}
