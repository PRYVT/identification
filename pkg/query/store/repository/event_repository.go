package repository

import (
	"database/sql"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	if db == nil {
		return nil
	}
	return &EventRepository{db: db}
}

func (repo *EventRepository) ReplaceEvent(newEventId string) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	// Delete old event
	stmt, err := tx.Prepare("DELETE FROM events")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Save new event
	stmt, err = tx.Prepare("INSERT INTO events (id) VALUES (?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newEventId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *EventRepository) GetLastEvent() (string, error) {
	var eventId string
	err := repo.db.QueryRow("SELECT id FROM events LIMIT 1").Scan(&eventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "0", nil
		}

		return "", err
	}
	return eventId, nil
}
