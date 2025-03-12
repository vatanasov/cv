package main

import (
	"autobiography/internal/database"
	"autobiography/internal/extractor"
	"autobiography/internal/models"
	"autobiography/internal/populator"
	"database/sql"
	"flag"
	"fmt"
	"os"
)

func main() {
	cvPath := flag.String("cv", "", "Path to cv PDF file")
	ghToken := flag.String("gh-token", "", "GitHub token")
	flag.Parse()

	db, err := database.New("./db.sqlite")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	if *cvPath != "" {
		populateDb(*cvPath, db)
	}

	if *ghToken != "" {
		getRepos(extractor.Token(*ghToken), db)
	}
}

func populateDb(cvPath string, db *sql.DB) {
	xml, err := extractor.ExtractXML(cvPath)
	if err != nil {
		fmt.Println("Error extracting XML:", err)
		os.Exit(1)
	}

	candidate, err := models.FromXML(xml)

	tx, err := db.Begin()
	m := models.NewModels(tx)
	err = populator.PopulateCandidate(m, &candidate)
	if err != nil {
		fmt.Println("Error populating candidate:", err)
		tx.Rollback()
		os.Exit(1)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		os.Exit(1)
	}
}

func getRepos(token extractor.Token, db *sql.DB) {
	repos, err := extractor.GitHubRepoApiClient.ExtractFromGitHub(token)
	if err != nil {
		fmt.Println("Error extracting GitHub repos:", err)
		os.Exit(1)
	}

	tx, err := db.Begin()
	m := models.NewModels(tx)
	err = populator.PopulateRepos(m, 1, repos)
	if err != nil {
		fmt.Println("Error populating candidate:", err)
		tx.Rollback()
		os.Exit(1)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		os.Exit(1)
	}
}
