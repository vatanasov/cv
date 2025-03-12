package models

import (
	"autobiography/internal/database"
)

type Communication struct {
	ChannelCode string `xml:"ChannelCode"`
	URI         string `xml:"URI"`
}

type CommunicationModel struct {
	DB database.Db
}

func (m *CommunicationModel) Insert(candidate Candidate) error {
	query := `
		INSERT INTO communications (candidate_id, channel_code, uri)
		VALUES (?, ?, ?)`

	for _, communicationChannel := range candidate.Communications {
		args := []any{candidate.ID, communicationChannel.ChannelCode, communicationChannel.URI}
		_, err := database.Insert(m.DB, query, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *CommunicationModel) GetAll(candidateId int64) ([]Communication, error) {
	query := `
		SELECT channel_code, uri
		FROM communications
		WHERE candidate_id = ?`

	args := []any{candidateId}

	communications, err := database.GetAll(
		m.DB,
		query,
		args,
		func(communication *Communication) []any {
			return []any{&communication.ChannelCode, &communication.URI}
		},
	)

	return communications, err
}
