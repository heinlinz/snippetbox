package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID        int
	Title     string
	Content   string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expireAt int) (int, error) {
	stmt := `INSERT INTO snippets(title, content, expires_at, created_at)
						VALUES(?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, content, expireAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The id returntype is int64, so we convert it into an int type before returning
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, expires_at, created_at FROM snippets
					WHERE expires_at > UTC_TIMESTAMP() and id = ?`
	rows := m.DB.QueryRow(stmt, id)

	// Initialize a nre zeroed snippet struct
	s := &Snippet{}

	// Use rows.scan to copy the data from each field
	err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, expires_at, created_at FROM snippets
					WHERE expires_at > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.ExpiresAt, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
